package cmds

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mitchell/selfpass/crypto"
)

func MakeGetCmd(masterpass string, cfg *viper.Viper, initClient CredentialClientInit) *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get a credential info and copy password to clipboard",
		Long: `Get a credential's info and copy password to clipboard, from Selfpass server, after
decrypting password.`,
		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			cbcontents, err := clipboard.ReadAll()
			check(err)

			restore := func(cbcontents string) {
				time.Sleep(time.Second * 5)
				clipboard.WriteAll(cbcontents)
			}

			cred, err := initClient(ctx).Get(ctx, args[0])
			check(err)

			key, err := hex.DecodeString(cfg.GetString(KeyPrivateKey))
			check(err)

			passkey, err := crypto.CombinePasswordAndKey([]byte(masterpass), key)
			check(err)

			passbytes, err := base64.StdEncoding.DecodeString(cred.Password)
			check(err)

			plainpass, err := crypto.CBCDecrypt(passkey, passbytes)

			check(clipboard.WriteAll(string(plainpass)))
			go restore(cbcontents)

			cjson, err := json.MarshalIndent(cred, "", "  ")
			check(err)
			fmt.Println(string(cjson))
			fmt.Println("Wrote password to clipboard.")
		},
	}

	return getCmd
}
