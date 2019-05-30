package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

func MakeList(initClient CredentialClientInit) *cobra.Command {
	var sourceHost string

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List the metadata for all credentials",
		Long: `List the metadata for all credentials, with the option to filter by source host. Metadata
includes almost all the information but the most sensitive.`,

		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()
			mdch, errch := initClient(ctx).GetAllMetadata(ctx, sourceHost)

			fmt.Println()

		receive:
			for count := 0; ; count++ {
				select {
				case <-ctx.Done():
					check(fmt.Errorf("context timeout"))

				case err := <-errch:
					check(err)

				case md, ok := <-mdch:
					if !ok {
						break receive
					}

					if count != 0 && count%3 == 0 {
						var proceed bool
						prompt := &survey.Confirm{Message: "Next page?", Default: true}
						check(survey.AskOne(prompt, &proceed, nil))

						if !proceed {
							break receive
						}

						clearCmd := exec.Command("clear")
						clearCmd.Stdout = os.Stdout
						check(clearCmd.Run())
					}

					fmt.Println(md)
				}
			}

			fmt.Println("Done listing.")
		},
	}

	listCmd.Flags().StringVarP(
		&sourceHost,
		"source-host",
		"s",
		"",
		"specify which source host to filter the results by",
	)

	return listCmd
}
