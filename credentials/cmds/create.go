package cmds

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/mitchell/selfpass/credentials/types"
	"github.com/mitchell/selfpass/crypto"
)

func MakeCreateCmd(masterpass string, cfg *viper.Viper, initClient CredentialClientInit) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a credential in Selfpass",
		Long: `Create a credential in Selfpass, and save it to the server after encrypting the
password.`,

		Run: func(_ *cobra.Command, args []string) {
			mdqs := []*survey.Question{
				{
					Name:   "primary",
					Prompt: &survey.Input{Message: "Primary user key:"},
				},
				{
					Name:   "sourceHost",
					Prompt: &survey.Input{Message: "Source host:"},
				},
				{
					Name:   "loginURL",
					Prompt: &survey.Input{Message: "Login url:"},
				},
				{
					Name:   "tag",
					Prompt: &survey.Input{Message: "Tag:"},
				},
			}
			cqs := []*survey.Question{
				{
					Name:   "username",
					Prompt: &survey.Input{Message: "Username:"},
				},
				{
					Name:   "email",
					Prompt: &survey.Input{Message: "Email:"},
				},
			}
			var ci types.CredentialInput

			check(survey.Ask(mdqs, &ci.MetadataInput))
			check(survey.Ask(cqs, &ci))

			var newpass bool
			prompt := &survey.Confirm{Message: "Do you want a random password?", Default: true}
			check(survey.AskOne(prompt, &newpass, nil))

			if newpass {
				ci.Password = generatePassword(16, true, true)

				var copypass bool
				prompt = &survey.Confirm{Message: "Copy new pass to clipboard?", Default: true}
				check(survey.AskOne(prompt, &copypass, nil))

				if copypass {
					check(clipboard.WriteAll(ci.Password))
				}
			} else {
				prompt := &survey.Password{Message: "Password: "}
				check(survey.AskOne(prompt, &ci.Password, nil))

				var cpass string
				prompt = &survey.Password{Message: "Confirm assword: "}
				check(survey.AskOne(prompt, &cpass, nil))

				if ci.Password != cpass {
					fmt.Println("passwords didn't match'")
					os.Exit(1)
				}
			}

			key, err := hex.DecodeString(cfg.GetString(KeyPrivateKey))
			check(err)

			keypass, err := crypto.CombinePasswordAndKey([]byte(masterpass), []byte(key))
			check(err)

			cipherpass, err := crypto.CBCEncrypt(keypass, []byte(ci.Password))
			check(err)

			ci.Password = base64.StdEncoding.EncodeToString(cipherpass)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			c, err := initClient(ctx).Create(ctx, ci)
			check(err)

			mdjson, err := json.MarshalIndent(c.Metadata, "", "  ")
			check(err)
			fmt.Println(string(mdjson))
		},
	}

	return createCmd
}

const alphas = "abcdefghijklmnopqrstuvABCDEFGHIJKLMNOPQRSTUV"
const alphanumerics = "abcdefghijklmnopqrstuvABCDEFGHIJKLMNOPQRSTUV1234567890"
const alphasAndSpecials = "abcdefghijklmnopqrstuvABCDEFGHIJKLMNOPQRSTUV1234567890!@#$%^&*()"
const alphanumericsAndSpecials = "abcdefghijklmnopqrstuvABCDEFGHIJKLMNOPQRSTUV1234567890!@#$%^&*()"

func generatePassword(length int, numbers, specials bool) string {
	rand.Seed(time.Now().UnixNano())
	pass := make([]byte, length)

	switch {
	case numbers && specials:
		for idx := 0; idx < length; idx++ {
			pass[idx] = alphanumericsAndSpecials[rand.Int63()%int64(len(alphanumericsAndSpecials))]
		}
	case numbers:
		for idx := 0; idx < length; idx++ {
			pass[idx] = alphanumerics[rand.Int63()%int64(len(alphanumerics))]
		}
	case specials:
		for idx := 0; idx < length; idx++ {
			pass[idx] = alphasAndSpecials[rand.Int63()%int64(len(alphasAndSpecials))]
		}
	default:
		for idx := 0; idx < length; idx++ {
			pass[idx] = alphas[rand.Int63()%int64(len(alphas))]
		}
	}

	return string(pass)
}
