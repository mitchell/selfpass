package commands

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"

	clitypes "github.com/mitchell/selfpass/cli/types"
	"github.com/mitchell/selfpass/credentials/types"
	"github.com/mitchell/selfpass/crypto"
)

func MakeCreate(repo clitypes.ConfigRepo, initClient CredentialClientInit) *cobra.Command {
	var length uint
	var numbers bool
	var specials bool

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a credential in Selfpass",
		Long: `Create a credential in Selfpass, and save it to the server after encrypting the
password.`,

		Run: func(_ *cobra.Command, args []string) {
			masterpass, cfg, err := repo.OpenConfig()
			check(err)

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

			key, err := hex.DecodeString(cfg.GetString(clitypes.KeyPrivateKey))
			check(err)

			keypass, err := crypto.CombinePasswordAndKey([]byte(masterpass), []byte(key))
			check(err)

			var newpass bool
			prompt := &survey.Confirm{Message: "Do you want a random password?", Default: true}
			check(survey.AskOne(prompt, &newpass, nil))

			if newpass {
				ci.Password = crypto.GeneratePassword(int(length), numbers, specials)

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
				prompt = &survey.Password{Message: "Confirm password: "}
				check(survey.AskOne(prompt, &cpass, nil))

				if ci.Password != cpass {
					fmt.Println("passwords didn't match")
					os.Exit(1)
				}
			}

			cipherpass, err := crypto.CBCEncrypt(keypass, []byte(ci.Password))
			check(err)

			ci.Password = base64.StdEncoding.EncodeToString(cipherpass)

			var otp bool
			prompt = &survey.Confirm{Message: "Do you have an OTP/MFA secret?", Default: true}
			check(survey.AskOne(prompt, &otp, nil))

			if otp {
				var secret string
				prompt := &survey.Password{Message: "OTP secret:"}
				check(survey.AskOne(prompt, &secret, nil))

				ciphersecret, err := crypto.CBCEncrypt(keypass, []byte(secret))
				check(err)

				ci.OTPSecret = base64.StdEncoding.EncodeToString(ciphersecret)

				var copyotp bool
				prompt2 := &survey.Confirm{Message: "Copy new OTP to clipboard?", Default: true}
				check(survey.AskOne(prompt2, &copyotp, nil))

				if copyotp {
					otp, err := totp.GenerateCode(secret, time.Now())
					check(err)

					check(clipboard.WriteAll(otp))
				}
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			c, err := initClient(ctx).Create(ctx, ci)
			check(err)

			fmt.Println(c)

			var cleancb bool
			prompt = &survey.Confirm{Message: "Do you want to clear the clipboard?", Default: true}
			check(survey.AskOne(prompt, &cleancb, nil))

			if cleancb {
				check(clipboard.WriteAll(" "))
			}
		},
	}

	createCmd.Flags().BoolVarP(&numbers, "numbers", "n", true, "use numbers in the generated password")
	createCmd.Flags().BoolVarP(&specials, "specials", "s", false, "use special characters in the generated password")
	createCmd.Flags().UintVarP(&length, "length", "l", 32, "length of the generated password")

	return createCmd
}
