package commands

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

func MakeDelete(initConfig CredentialClientInit) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a credential using the given ID",
		Long:  `Delete a credential using the given ID, permanently. THERE IS NO UNDOING THIS ACTION.`,
		Args:  cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			var confirmed bool
			prompt := &survey.Confirm{Message: "Are you sure you want to permanently delete this credential?"}
			check(survey.AskOne(prompt, &confirmed, nil))

			if confirmed {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				check(initConfig(ctx).Delete(ctx, args[0]))
			}
		},
	}

	return deleteCmd
}
