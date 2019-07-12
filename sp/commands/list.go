package commands

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/mitchell/selfpass/services/credentials/types"
)

func makeList(initClient CredentialsClientInit) *cobra.Command {
	var sourceHost string

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List the metadata for all credentials",
		Long: `List the metadata for all credentials, with the option to filter by source host. Metadata
includes almost all the information but the most sensitive.`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
			defer cancel()

			mdch, errch := initClient(ctx).GetAllMetadata(ctx, sourceHost)
			mds := map[string][]types.Metadata{}

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

					mds[md.SourceHost] = append(mds[md.SourceHost], md)
				}
			}

			sources := []string{}
			for source := range mds {
				sources = append(sources, source)
			}

			sort.Strings(sources)

			prompt := &survey.Select{
				Message:  "Source host:",
				Options:  sources,
				PageSize: 20,
				VimMode:  true,
			}

			var source string
			check(survey.AskOne(prompt, &source, nil))

			for _, md := range mds[source] {
				fmt.Println(md)
			}

			fmt.Println("Done listing.")
		},
	}

	return listCmd
}
