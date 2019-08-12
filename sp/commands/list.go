package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/mitchell/selfpass/services/credentials/types"
)

func makeList(initClient credentialsClientInit) *cobra.Command {
	flags := credentialFlagSet{}.withHostFlag()

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List the metadata for all credentials",
		Long: `List the metadata for all credentials, with the option to filter by source host. Metadata
includes almost all the information but the most sensitive.`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
			defer cancel()

			client := initClient(ctx)

			ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			mdch, errch := client.GetAllMetadata(ctx, flags.sourceHost)
			var mds []types.Metadata

			fmt.Println()

		receive:
			for {
				select {
				case <-ctx.Done():
					check(ctx.Err())

				case err := <-errch:
					check(err)

				case md, ok := <-mdch:
					if !ok {
						break receive
					}

					mds = append(mds, md)
				}
			}

			var sources []string
			mdmap := map[string][]types.Metadata{}
			for _, md := range mds {
				tmds := mdmap[md.SourceHost]

				if tmds == nil {
					mdmap[md.SourceHost] = []types.Metadata{md}
					sources = append(sources, md.SourceHost)
					continue
				}

				mdmap[md.SourceHost] = append(mdmap[md.SourceHost], md)
			}

			if flags.sourceHost == "" {
				prompt := &survey.Select{
					Message:  "Source host:",
					Options:  sources,
					PageSize: 20,
					VimMode:  true,
				}

				check(survey.AskOne(prompt, &flags.sourceHost, nil))
			}

			if len(mdmap[flags.sourceHost]) == 0 {
				check(errSourceNotFound)
			}

			for _, md := range mdmap[flags.sourceHost] {
				fmt.Println(md)
			}

			fmt.Println("Done listing.")
		},
	}

	flags.register(listCmd)

	return listCmd
}
