package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/mitchell/selfpass/credentials/cmds"
	"github.com/mitchell/selfpass/credentials/types"
	"github.com/mitchell/selfpass/crypto"
)

func Execute(ctx context.Context, initClient types.CredentialClientInit) {
	rootCmd := &cobra.Command{
		Use:   "spc",
		Short: "This is the CLI client for Selfpass.",
		Long: `This is the CLI client for Selfpass, the self-hosted password manager. With this tool you
can interact with the entire Selfpass API.`,
	}

	cfgFile := rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.spc.toml)")
	rootCmd.PersistentFlags().Parse(os.Args)

	decryptCfg := rootCmd.Flags().Bool("decrypt-cfg", false, "unencrypt config file")
	rootCmd.Flags().Parse(os.Args)

	encryptCfg := !*decryptCfg
	masterpass, cfg := openConfig(*cfgFile)
	if encryptCfg && masterpass != "" {
		defer encryptConfig(masterpass, cfg)
	}
	if *decryptCfg {
		fmt.Println("Decrypting config file. It will auto-encrypt when you next run of spc.")
		return
	}

	rootCmd.AddCommand(makeInitCmd(cfg))
	rootCmd.AddCommand(cmds.MakeListCmd(makeInitClient(cfg, initClient)))
	rootCmd.AddCommand(cmds.MakeCreateCmd(masterpass, cfg, makeInitClient(cfg, initClient)))
	rootCmd.AddCommand(cmds.MakeGetCmd(masterpass, cfg, makeInitClient(cfg, initClient)))

	check(rootCmd.Execute())
}

func makeInitClient(cfg *viper.Viper, initClient types.CredentialClientInit) cmds.CredentialClientInit {
	return func(ctx context.Context) types.CredentialClient {
		connConfig := cfg.GetStringMapString(cmds.KeyConnConfig)

		client, err := initClient(
			ctx,
			connConfig["target"],
			connConfig["ca"],
			connConfig["cert"],
			connConfig["key"],
		)
		if err != nil {
			fmt.Printf("Please run 'init' command before running API commands.\nError Message: %s\n", err)
			os.Exit(1)
		}

		return client
	}
}

func openConfig(cfgFile string) (masterpass string, v *viper.Viper) {
	v = viper.New()
	v.SetConfigType("toml")

	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		check(err)

		// Search config in home directory with name ".spc" (without extension).
		v.AddConfigPath(home)
		v.SetConfigName(".spc")

		cfgFile = home + "/.spc.toml"
	}

	if _, err := os.Open(cfgFile); !os.IsNotExist(err) {
		prompt := &survey.Password{Message: "Master password:"}
		check(survey.AskOne(prompt, &masterpass, nil))

		decryptConfig(masterpass, cfgFile)
	}

	//v.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", v.ConfigFileUsed())
	}

	return masterpass, v
}

func decryptConfig(masterpass string, cfgFile string) {
	contents, err := ioutil.ReadFile(cfgFile)
	check(err)

	passkey, err := crypto.GenerateKeyFromPassword([]byte(masterpass))
	check(err)

	contents, err = crypto.CBCDecrypt(passkey, contents)
	if err != nil && err.Error() == "Padding incorrect" {
		fmt.Println("incorrect master password")
		os.Exit(1)
	} else if err != nil && err.Error() == "ciphertext is not a multiple of the block size" {
		fmt.Println("Config wasn't encrypted.")
		return
	}
	check(err)

	check(ioutil.WriteFile(cfgFile, contents, 0600))
}

func encryptConfig(masterpass string, cfg *viper.Viper) {
	contents, err := ioutil.ReadFile(cfg.ConfigFileUsed())
	if os.IsNotExist(err) {
		return
	}

	keypass, err := crypto.GenerateKeyFromPassword([]byte(masterpass))
	check(err)

	contents, err = crypto.CBCEncrypt(keypass, contents)
	check(err)

	check(ioutil.WriteFile(cfg.ConfigFileUsed(), contents, 0600))
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
