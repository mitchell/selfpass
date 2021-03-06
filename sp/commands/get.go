package commands

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/mitchell/selfpass/sp/crypto"
	clitypes "github.com/mitchell/selfpass/sp/types"
)

func makeGet(repo clitypes.ConfigRepo, initClient credentialsClientInit) *cobra.Command {
	flags := credentialFlagSet{}.withCredFlags()

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a credential info and copy password to clipboard",
		Long: `Get a credential's info and copy password to clipboard, from Selfpass server, after
decrypting password.`,

		Run: func(cmd *cobra.Command, args []string) {
			defer flags.resetValues()

			var (
				copyPass bool
				cleancb  bool
				prompt   survey.Prompt
			)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
			defer cancel()

			client := initClient(ctx)
			masterpass, cfg, err := repo.OpenConfig()
			check(err)

			cred := selectCredential(client, flags.sourceHost, flags.primary)

			fmt.Println(cred)

			check(clipboard.WriteAll(string(cred.Primary)))

			fmt.Println("Wrote primary user key to clipboard.")

			key := cfg.GetString(clitypes.KeyPrivateKey)
			passkey := crypto.GeneratePBKDF2Key([]byte(masterpass), []byte(key))

			prompt = &survey.Confirm{Message: "Copy password to clipboard?", Default: true}
			check(survey.AskOne(prompt, &copyPass, nil))

			if copyPass {
				passbytes, err := base64.StdEncoding.DecodeString(cred.Password)
				check(err)

				plainpass, err := crypto.CBCDecrypt(passkey, passbytes)

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

					plainsecret, err := crypto.CBCDecrypt(passkey, secretbytes)

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
				fmt.Println("Clipboard cleared.")
			}
		},
	}

	flags.register(getCmd)

	return getCmd
}
