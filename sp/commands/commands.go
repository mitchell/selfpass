package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/mitchell/selfpass/services/credentials/types"
)

type credentialsClientInit func(ctx context.Context) (c types.CredentialsClient)

var errSourceNotFound = errors.New("source host not found")

type credentialFlagSet struct {
	includePasswordFlags bool
	includeHostFlag      bool

	sourceHost string
	noNumbers  bool
	noSpecials bool
	length     uint
}

func (set credentialFlagSet) withPasswordFlags() credentialFlagSet {
	set.includePasswordFlags = true
	return set
}

func (set credentialFlagSet) withHostFlag() credentialFlagSet {
	set.includeHostFlag = true
	return set
}

func (set *credentialFlagSet) register(cmd *cobra.Command) {
	if set.includeHostFlag {
		cmd.Flags().StringVarP(&set.sourceHost, "source-host", "s", "", "filter results to this source host")
	}

	if set.includePasswordFlags {
		cmd.Flags().BoolVarP(&set.noNumbers, "no-numbers", "n", false, "do not use numbers in the generated password")
		cmd.Flags().BoolVarP(&set.noSpecials, "no-specials", "p", false, "do not use special characters in the generated password")
		cmd.Flags().UintVarP(&set.length, "length", "l", 32, "length of the generated password")
	}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func selectCredential(client types.CredentialsClient, sourceHost string) types.Credential {
	var (
		idKey  string
		prompt survey.Prompt
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mdch, errch := client.GetAllMetadata(ctx, sourceHost)
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

	if sourceHost == "" {
		prompt = &survey.Select{
			Message:  "Source host:",
			Options:  sources,
			PageSize: 20,
			VimMode:  true,
		}

		check(survey.AskOne(prompt, &sourceHost, nil))
	}

	if len(mdmap[sourceHost]) == 0 {
		check(errSourceNotFound)
	}

	keys := []string{}
	keyIDMap := map[string]string{}
	for _, md := range mdmap[sourceHost] {
		key := md.Primary
		if md.Tag != "" {
			key += "-" + md.Tag
		}
		keys = append(keys, key)
		keyIDMap[key] = md.ID
	}

	prompt = &survey.Select{
		Message:  "Primary user key (and tag):",
		Options:  keys,
		PageSize: 20,
		VimMode:  true,
	}

	check(survey.AskOne(prompt, &idKey, nil))

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cred, err := client.Get(ctx, keyIDMap[idKey])
	check(err)

	return cred
}
