package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

func makeDelete(initClient CredentialClientInit) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a credential using the given ID",
		Long:  `Delete a credential using the given ID, permanently. THERE IS NO UNDOING THIS ACTION.`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
			defer cancel()

			client := initClient(ctx)

			cred := selectCredential(client)

			fmt.Println(cred)

			var confirmed bool
			prompt := &survey.Confirm{Message: "Are you sure you want to permanently delete this credential?"}
			check(survey.AskOne(prompt, &confirmed, nil))

			if confirmed {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*25)
				defer cancel()

				check(initClient(ctx).Delete(ctx, cred.ID))
			}
		},
	}

	return deleteCmd
}
