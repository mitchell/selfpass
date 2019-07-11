package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mitchell/selfpass/services/cli/repositories"
	"github.com/mitchell/selfpass/services/cli/types"
	"github.com/mitchell/selfpass/services/credentials/commands"
	credrepos "github.com/mitchell/selfpass/services/credentials/repositories"
	credtypes "github.com/mitchell/selfpass/services/credentials/types"
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
	clientInit := credrepos.NewCredentialServiceClient

	rootCmd.AddCommand(
		makeInit(mgr),
		makeEncrypt(mgr),
		makeDecrypt(mgr),
		makeDecryptCfg(mgr),
		commands.MakeList(makeInitClient(mgr, clientInit)),
		commands.MakeCreate(mgr, makeInitClient(mgr, clientInit)),
		commands.MakeUpdate(mgr, makeInitClient(mgr, clientInit)),
		commands.MakeGet(mgr, makeInitClient(mgr, clientInit)),
		commands.MakeDelete(makeInitClient(mgr, clientInit)),
		commands.MakeGCMToCBC(mgr, makeInitClient(mgr, clientInit)),
	)

	check(rootCmd.Execute())
}

func makeInitClient(repo types.ConfigRepo, initClient credtypes.CredentialClientInit) commands.CredentialClientInit {
	return func(ctx context.Context) credtypes.CredentialClient {
		_, cfg, err := repo.OpenConfig()
		check(err)

		connConfig := cfg.GetStringMapString(types.KeyConnConfig)

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
