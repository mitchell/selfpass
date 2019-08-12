package commands

import (
	"context"

	"github.com/spf13/cobra"

	credrepos "github.com/mitchell/selfpass/services/credentials/repositories"
	credtypes "github.com/mitchell/selfpass/services/credentials/types"
	"github.com/mitchell/selfpass/sp/repositories"
	"github.com/mitchell/selfpass/sp/types"
)

// Execute is the main entrypoint for the `sp` CLI tool.
func Execute() {
	rootCmd := &cobra.Command{
		Use:   "sp",
		Short: "This is the CLI client for Selfpass.",
		Long: `This is the CLI client for Selfpass, the self-hosted password manager. With this tool you
can interact with the entire Selfpass API.`,
		Version: "v0.1.0",
	}

	cfgFile := rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.sp.toml)")

	mgr := repositories.NewConfigManager(cfgFile)
	clientInit := credrepos.NewCredentialsClient

	rootCmd.AddCommand(
		makeInit(mgr),
		makeEncrypt(mgr),
		makeDecrypt(mgr),
		makeDecryptCfg(mgr),
		makeList(makeInitClient(mgr, clientInit)),
		makeCreate(mgr, makeInitClient(mgr, clientInit)),
		makeUpdate(mgr, makeInitClient(mgr, clientInit)),
		makeGet(mgr, makeInitClient(mgr, clientInit)),
		makeDelete(makeInitClient(mgr, clientInit)),
		makeGCMToCBC(mgr, makeInitClient(mgr, clientInit)),
	)

	check(rootCmd.Execute())
}

func makeInitClient(repo types.ConfigRepo, initClient credtypes.CredentialsClientInit) credentialsClientInit {
	return func(ctx context.Context) credtypes.CredentialsClient {
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
