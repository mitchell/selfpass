package commands

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/mitchell/selfpass/sp/types"
)

func makeInit(repo types.ConfigRepo) *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "This command initializes SPC for the first time",
		Long: `This command initializes SPC for the first time. Writing to the user configuration
the users private key, and server certificates. (All of which will be encrypted)`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				hasPK       bool
				masterpass  string
				cmasterpass string
				target      string
				caFile      string
				certFile    string
				keyFile     string
				prompt      survey.Prompt
				privateKey  = strings.Replace(uuid.New().String(), "-", "", -1)
			)
			_, cfg, _ := repo.OpenConfig()

			prompt = &survey.Password{Message: "New master password:"}
			check(survey.AskOne(prompt, &masterpass, nil))

			prompt = &survey.Password{Message: "Confirm master password:"}
			check(survey.AskOne(prompt, &cmasterpass, nil))
			if masterpass != cmasterpass {
				check(fmt.Errorf("master passwords didn't match"))
			}

			repo.SetMasterpass(masterpass)

			prompt = &survey.Input{Message: "Selfpass server address:"}
			check(survey.AskOne(prompt, &target, nil))

			prompt = &survey.Confirm{Message: "Do you have a private key?"}
			check(survey.AskOne(prompt, &hasPK, nil))

			if hasPK {
				prompt = &survey.Password{Message: "Private key:"}
				check(survey.AskOne(prompt, &privateKey, nil))
				privateKey = strings.Replace(privateKey, "-", "", -1)
			}

			prompt = &survey.Input{Message: "CA certificate file:"}
			check(survey.AskOne(prompt, &caFile, nil))
			ca, err := ioutil.ReadFile(caFile)
			check(err)

			prompt = &survey.Input{Message: "Client certificate file:"}
			check(survey.AskOne(prompt, &certFile, nil))
			cert, err := ioutil.ReadFile(certFile)
			check(err)

			prompt = &survey.Input{Message: "Client key file:"}
			check(survey.AskOne(prompt, &keyFile, nil))
			key, err := ioutil.ReadFile(keyFile)
			check(err)

			cfg.Set(types.KeyConnConfig, map[string]string{
				"target": target,
				"ca":     string(ca),
				"cert":   string(cert),
				"key":    string(key),
			})

			cfg.Set(types.KeyPrivateKey, privateKey)

			check(repo.WriteConfig())
		},
	}

	return initCmd
}
