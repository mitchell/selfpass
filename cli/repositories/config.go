package repositories

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/mitchell/selfpass/crypto"
)

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

		cfg = home + "/.spc.toml"
	}

	mgr.v.SetConfigFile(cfg)
	mgr.cfgFile = &cfg

	var contents []byte
	var cipherAuthFailed bool

	if _, err := os.Open(cfg); os.IsNotExist(err) {
		return output, mgr.v, fmt.Errorf("no config found, run 'init' command")
	}

	prompt := &survey.Password{Message: "Master password:"}
	if err = survey.AskOne(prompt, &mgr.masterpass, nil); err != nil {
		return output, nil, err
	}

	contents, err = decryptConfig(mgr.masterpass, cfg)
	if err != nil && err.Error() == "cipher: message authentication failed" {
		cipherAuthFailed = true
	} else if err != nil {
		return output, nil, err
	}

	if err = mgr.v.ReadConfig(bytes.NewBuffer(contents)); err != nil && err.Error() == "While parsing config: (1, 1): unexpected token" {
		return output, nil, fmt.Errorf("incorrect master password")
	} else if err != nil {
		return output, nil, err
	}

	if cipherAuthFailed {
		fmt.Println("Config wasn't encrypted, or has been compromised.")

		if err = mgr.WriteConfig(); err != nil {
			return output, nil, err
		}
	}

	return mgr.masterpass, mgr.v, nil
}

func decryptConfig(masterpass string, cfgFile string) (contents []byte, err error) {
	contents, err = ioutil.ReadFile(cfgFile)
	if err != nil {
		return contents, err
	}

	passkey, err := crypto.GenerateKeyFromPassword([]byte(masterpass))
	if err != nil {
		return contents, err
	}

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

	keypass, err := crypto.GenerateKeyFromPassword([]byte(mgr.masterpass))
	if err != nil {
		return err
	}

	contents, err = crypto.GCMEncrypt(keypass, contents)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(mgr.v.ConfigFileUsed(), contents, 0600)
	if err != nil {
		return err
	}

	return nil
}
