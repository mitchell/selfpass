package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func MakeListCmd(initClient CredentialClientInit) *cobra.Command {
	var sourceHost string

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List the metadata for all credentials",
		Long: `List the metadata for all credentials, with the option to filter by source host. Metadata
includes almost all the information but the most sensitive.`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			mdch, errch := initClient(ctx).GetAllMetadata(ctx, sourceHost)

		receive:
			for {
				select {
				case <-ctx.Done():
					check(fmt.Errorf("context timeout"))

				case err := <-errch:
					check(err)

				case md, ok := <-mdch:
					if !ok {
						break receive
					}

					mdjson, err := json.MarshalIndent(md, "", "  ")
					check(err)
					fmt.Println(string(mdjson))
				}
			}
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
