package repositories

import (
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
	decrypted  bool
	decrypt    bool
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

	v = viper.New()
	mgr.v = v

	v.SetConfigType("toml")

	if cfg != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfg)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return output, nil, err
		}

		// Search config in home directory with name ".spc" (without extension).
		v.AddConfigPath(home)
		v.SetConfigName(".spc")

		cfg = home + "/.spc.toml"
	}

	if _, err := os.Open(cfg); !os.IsNotExist(err) {
		prompt := &survey.Password{Message: "Master password:"}
		if err = survey.AskOne(prompt, &mgr.masterpass, nil); err != nil {
			return output, nil, err
		}

		mgr.decrypted, err = decryptConfig(mgr.masterpass, cfg)
		if err != nil {
			return output, nil, err
		}
	}

	// v.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err = v.ReadInConfig(); err != nil {
		mgr.decrypted = true
		return output, v, fmt.Errorf("no config found, run 'init' command")
	}

	return mgr.masterpass, mgr.v, nil
}

func decryptConfig(masterpass string, cfgFile string) (decrypted bool, err error) {
	contents, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return decrypted, err
	}

	passkey, err := crypto.GenerateKeyFromPassword([]byte(masterpass))
	if err != nil {
		return decrypted, err
	}

	contents, err = crypto.CBCDecrypt(passkey, contents)
	if err != nil && err.Error() == "Padding incorrect" {
		return decrypted, fmt.Errorf("incorrect master password")
	} else if err != nil && err.Error() == "ciphertext is not a multiple of the block size" {
		fmt.Println("Config wasn't encrypted.")
		return true, nil
	}
	if err != nil {
		return decrypted, err
	}

	if err = ioutil.WriteFile(cfgFile, contents, 0600); err != nil {
		return decrypted, err
	}

	return true, nil
}

func (mgr *ConfigManager) DecryptConfig() {
	mgr.decrypt = true
}

func (mgr *ConfigManager) CloseConfig() {
	if !mgr.decrypt && mgr.decrypted {
		contents, err := ioutil.ReadFile(mgr.v.ConfigFileUsed())
		if os.IsNotExist(err) {
			return
		}

		keypass, err := crypto.GenerateKeyFromPassword([]byte(mgr.masterpass))
		if err != nil {
			panic(err)
		}

		contents, err = crypto.CBCEncrypt(keypass, contents)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(mgr.v.ConfigFileUsed(), contents, 0600)
		if err != nil {
			panic(err)
		}

		return
	}

	if mgr.decrypt {
		fmt.Println("Decrypting config file. It will auto-encrypt when you next run spc.")
	}
}
