package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mitchell/selfpass/cli/repositories"
	"github.com/mitchell/selfpass/cli/types"
	"github.com/mitchell/selfpass/credentials/commands"
	credrepos "github.com/mitchell/selfpass/credentials/repositories"
	credtypes "github.com/mitchell/selfpass/credentials/types"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "spc",
		Short: "This is the CLI client for Selfpass.",
		Long: `This is the CLI client for Selfpass, the self-hosted password manager. With this tool you
can interact with the entire Selfpass API.`,
		Version: "v0.1.0",
	}

	cfgFile := rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.spc.toml)")

	mgr := repositories.NewConfigManager(cfgFile)
	defer mgr.CloseConfig()

	clientInit := credrepos.NewCredentialServiceClient

	rootCmd.AddCommand(makeInit(mgr))
	rootCmd.AddCommand(makeEncrypt(mgr))
	rootCmd.AddCommand(makeDecrypt(mgr))
	rootCmd.AddCommand(makeDecryptCfg(mgr))
	rootCmd.AddCommand(commands.MakeList(makeInitClient(mgr, clientInit)))
	rootCmd.AddCommand(commands.MakeCreate(mgr, makeInitClient(mgr, clientInit)))
	rootCmd.AddCommand(commands.MakeGet(mgr, makeInitClient(mgr, clientInit)))
	rootCmd.AddCommand(commands.MakeDelete(makeInitClient(mgr, clientInit)))

	check(rootCmd.Execute())
}

func makeInitClient(repo types.ConfigRepo, initClient credtypes.CredentialClientInit) commands.CredentialClientInit {
	return func(ctx context.Context) credtypes.CredentialClient {
		_, cfg, err := repo.OpenConfig()
		check(err)

		connConfig := cfg.GetStringMapString(keyConnConfig)

		client, err := initClient(
			ctx,
			connConfig["target"],
			connConfig["ca"],
			connConfig["cert"],
			connConfig["key"],
		)
		check(err)

		return client
	}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const keyConnConfig = "connection"
