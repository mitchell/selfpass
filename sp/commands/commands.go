package commands

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/mitchell/selfpass/services/credentials/types"
)

type CredentialsClientInit func(ctx context.Context) (c types.CredentialsClient)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func selectCredential(client types.CredentialsClient) types.Credential {
	var (
		idKey  string
		source string
		prompt survey.Prompt
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mdch, errch := client.GetAllMetadata(ctx, "")
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

	prompt = &survey.Select{
		Message:  "Source host:",
		Options:  sources,
		PageSize: 20,
		VimMode:  true,
	}

	check(survey.AskOne(prompt, &source, nil))

	keys := []string{}
	keyIDMap := map[string]string{}
	for _, md := range mds[source] {
		key := md.Primary
		if md.Tag != "" {
			key += "-" + md.Tag
		}
		keys = append(keys, key)
		keyIDMap[key] = md.ID
	}

	sort.Strings(keys)

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
