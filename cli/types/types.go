package types

import "github.com/spf13/viper"

type ConfigRepo interface {
	OpenConfig() (masterpass string, cfg *viper.Viper, err error)
	DecryptConfig()
	SetMasterpass(masterpass string)
}
