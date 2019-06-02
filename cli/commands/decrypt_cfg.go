package commands

import (
	"github.com/spf13/cobra"

	"github.com/mitchell/selfpass/cli/types"
)

func makeDecryptCfg(repo types.ConfigRepo) *cobra.Command {
	decryptCfg := &cobra.Command{
		Use:   "decrypt-cfg",
		Short: "Decrypt your config file",
		Long:  "Decrypt your config file, so you can access your private key, host, and certs.",

		Run: func(cmd *cobra.Command, args []string) {
			_, _, err := repo.OpenConfig()
			check(err)

			repo.DecryptConfig()
		},
	}

	return decryptCfg
}
