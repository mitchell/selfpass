package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
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
		Run:   run,
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

func run(cmd *cobra.Command, _ []string) {
	ss := []prompt.Suggest{
		{Text: "exit", Description: "Exit selfpass prompt"},
	}

	for _, subcmd := range cmd.Commands() {
		ss = append(ss, prompt.Suggest{
			Text:        subcmd.Name(),
			Description: subcmd.Short,
		})
	}

	completer := func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(ss, d.TextBeforeCursor(), true)
	}

	executor := func(argstr string) {
		args := strings.Split(argstr, " ")

		if len(args) > 0 && args[0] == "exit" {
			fmt.Print("Goodbye!\n\n")
			os.Exit(0)
		}

		cmd.SetArgs(args)
		err := cmd.Execute()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println()
	}

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("selfpass-> "),
		prompt.OptionPrefixTextColor(prompt.White),
		prompt.OptionSuggestionBGColor(prompt.DarkBlue),
		prompt.OptionDescriptionBGColor(prompt.Blue),
	)

	p.Run()
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
