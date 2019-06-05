package types

import "github.com/spf13/viper"

type ConfigRepo interface {
	OpenConfig() (masterpass string, cfg *viper.Viper, err error)
	DecryptConfig() (err error)
	SetMasterpass(masterpass string)
	WriteConfig() (err error)
}

const KeyPrivateKey = "private_key"
const KeyConnConfig = "connection"
