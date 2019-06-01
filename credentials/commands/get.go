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
	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/mitchell/selfpass/credentials/types"
	"github.com/mitchell/selfpass/crypto"
)

func MakeGet(masterpass string, cfg *viper.Viper, initClient CredentialClientInit) *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a credential info and copy password to clipboard",
		Long: `Get a credential's info and copy password to clipboard, from Selfpass server, after
decrypting password.`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			client := initClient(ctx)
			mdch, errch := client.GetAllMetadata(ctx, "")
			mds := map[string][]types.Metadata{}

			fmt.Println()

		receive:
			for count := 0; ; count++ {
				select {
				case <-ctx.Done():
					check(fmt.Errorf("context timeout"))

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

			var prompt survey.Prompt
			prompt = &survey.Select{
				Message:  "Source host:",
				Options:  sources,
				PageSize: 10,
				VimMode:  true,
			}

			var source string
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

			prompt = &survey.Select{
				Message:  "Primary user key (and tag):",
				Options:  keys,
				PageSize: 10,
				VimMode:  true,
			}

			var idKey string
			check(survey.AskOne(prompt, &idKey, nil))

			ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			cred, err := client.Get(ctx, keyIDMap[idKey])
			check(err)

			fmt.Println(cred)

			check(clipboard.WriteAll(string(cred.Primary)))

			fmt.Println("Wrote primary user key to clipboard.")

			key, err := hex.DecodeString(cfg.GetString(KeyPrivateKey))
			check(err)

			passkey, err := crypto.CombinePasswordAndKey([]byte(masterpass), key)
			check(err)

			var copyPass bool
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

			var cleancb bool
			prompt = &survey.Confirm{Message: "Do you want to clear the clipboard?", Default: true}
			check(survey.AskOne(prompt, &cleancb, nil))

			if cleancb {
				check(clipboard.WriteAll(" "))
			}
		},
	}

	return getCmd
}
