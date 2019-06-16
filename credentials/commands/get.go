package commands

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"

	clitypes "github.com/mitchell/selfpass/cli/types"
	"github.com/mitchell/selfpass/crypto"
)

func MakeGet(repo clitypes.ConfigRepo, initClient CredentialClientInit) *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a credential info and copy password to clipboard",
		Long: `Get a credential's info and copy password to clipboard, from Selfpass server, after
decrypting password.`,

		Run: func(cmd *cobra.Command, args []string) {
			var (
				copyPass bool
				cleancb  bool
				prompt   survey.Prompt
			)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
			defer cancel()

			client := initClient(ctx)
			masterpass, cfg, err := repo.OpenConfig()
			check(err)

			cred := selectCredential(client)

			fmt.Println(cred)

			check(clipboard.WriteAll(string(cred.Primary)))

			fmt.Println("Wrote primary user key to clipboard.")

			key, err := hex.DecodeString(cfg.GetString(clitypes.KeyPrivateKey))
			check(err)

			passkey, err := crypto.CombinePasswordAndKey([]byte(masterpass), key)
			check(err)

			prompt = &survey.Confirm{Message: "Copy password to clipboard?", Default: true}
			check(survey.AskOne(prompt, &copyPass, nil))

			if copyPass {
				passbytes, err := base64.StdEncoding.DecodeString(cred.Password)
				check(err)

				plainpass, err := crypto.GCMDecrypt(passkey, passbytes)

				check(clipboard.WriteAll(string(plainpass)))

				fmt.Println("Wrote password to clipboard.")
			}

			if cred.OTPSecret != "" {
				var newOTP bool
				prompt = &survey.Confirm{Message: "Generate one time password and copy to clipboard?", Default: true}
				check(survey.AskOne(prompt, &newOTP, nil))

				if newOTP {
					secretbytes, err := base64.StdEncoding.DecodeString(cred.OTPSecret)
					check(err)

					plainsecret, err := crypto.GCMDecrypt(passkey, secretbytes)

					otp, err := totp.GenerateCode(string(plainsecret), time.Now())
					check(err)

					check(clipboard.WriteAll(otp))

					fmt.Println("Wrote one time password to clipboard.")
				}
			}

			prompt = &survey.Confirm{Message: "Do you want to clear the clipboard?", Default: true}
			check(survey.AskOne(prompt, &cleancb, nil))

			if cleancb {
				check(clipboard.WriteAll(" "))
			}
		},
	}

	return getCmd
}
