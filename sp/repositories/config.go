package repositories

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/mitchell/selfpass/sp/crypto"
	"github.com/mitchell/selfpass/sp/types"
)

var ErrNoConfigFound = errors.New("no config found, run 'init' command")

func NewConfigManager(cfgFile *string) *ConfigManager {
	return &ConfigManager{
		cfgFile: cfgFile,
	}
}

type ConfigManager struct {
	masterpass string
	cfgFile    *string
	v          *viper.Viper
}

func (mgr *ConfigManager) SetMasterpass(masterpass string) {
	mgr.masterpass = masterpass
}

func (mgr *ConfigManager) OpenConfig() (output string, v *viper.Viper, err error) {
	if mgr.masterpass != "" {
		return mgr.masterpass, mgr.v, nil
	}
	cfg := *mgr.cfgFile

	mgr.v = viper.New()
	mgr.v.SetConfigType("toml")

	if cfg == "" {
		home, err := homedir.Dir()
		if err != nil {
			return output, nil, err
		}

		cfg = home + "/.sp.toml"
	}

	mgr.v.SetConfigFile(cfg)
	mgr.cfgFile = &cfg

	var contents []byte
	var configDecrypted bool

	if _, err := os.Open(cfg); os.IsNotExist(err) {
		return output, mgr.v, ErrNoConfigFound
	}

	prompt := &survey.Password{Message: "Master password:"}
	if err = survey.AskOne(prompt, &mgr.masterpass, nil); err != nil {
		return output, nil, err
	}

	contents, err = decryptConfig(mgr.masterpass, cfg)
	if err == errConfigDecrypted {
		configDecrypted = true
	} else if err != nil && err.Error() == crypto.ErrAuthenticationFailed.Error() {
		return output, nil, errors.New("incorrect masterpass")
	} else if err != nil {
		return output, nil, err
	}

	if err = mgr.v.ReadConfig(bytes.NewBuffer(contents)); err != nil {
		return output, nil, err
	}

	if configDecrypted {
		fmt.Println("Config wasn't encrypted, or has been compromised.")

		if err = mgr.WriteConfig(); err != nil {
			return output, nil, err
		}
	}

	return mgr.masterpass, mgr.v, nil
}

var errConfigDecrypted = errors.New("config is decrypted")

func decryptConfig(masterpass string, cfgFile string) (contents []byte, err error) {
	contents, err = ioutil.ReadFile(cfgFile)
	if err != nil {
		return contents, err
	}

	if string(contents[:len(types.KeyPrivateKey)]) == types.KeyPrivateKey {
		return contents, errConfigDecrypted
	}

	salt := contents[:saltSize]
	contents = contents[saltSize:]

	passkey := crypto.GeneratePBKDF2Key([]byte(masterpass), salt)

	plaintext, err := crypto.GCMDecrypt(passkey, contents)
	if err != nil {
		return contents, err
	}

	return plaintext, nil
}

func (mgr ConfigManager) DecryptConfig() error {
	if err := mgr.v.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func (mgr ConfigManager) WriteConfig() (err error) {
	if err := mgr.v.WriteConfigAs(*mgr.cfgFile); err != nil {
		return err
	}

	contents, err := ioutil.ReadFile(mgr.v.ConfigFileUsed())
	if err != nil {
		return err
	}

	salt := make([]byte, saltSize)
	_, err = rand.Read(salt)
	if err != nil {
		return err
	}

	keypass := crypto.GeneratePBKDF2Key([]byte(mgr.masterpass), salt)

	contents, err = crypto.GCMEncrypt(keypass, contents)
	if err != nil {
		return err
	}

	contents = append(salt, contents...)

	err = ioutil.WriteFile(mgr.v.ConfigFileUsed(), contents, 0600)
	if err != nil {
		return err
	}

	return nil
}

const saltSize = 16
