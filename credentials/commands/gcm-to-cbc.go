package commands

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	clitypes "github.com/mitchell/selfpass/cli/types"
	"github.com/mitchell/selfpass/credentials/types"
	"github.com/mitchell/selfpass/crypto"
	"github.com/spf13/cobra"
)

func MakeGCMToCBC(repo clitypes.ConfigRepo, initClient CredentialClientInit) *cobra.Command {
	gcmToCBC := &cobra.Command{
		Use:    "gcm-to-cbc",
		Hidden: true,

		Run: func(cmd *cobra.Command, args []string) {
			masterpass, cfg, err := repo.OpenConfig()
			check(err)

			privKey := cfg.GetString(clitypes.KeyPrivateKey)

			fmt.Println(privKey)

			oldHex, err := hex.DecodeString(privKey)
			check(err)

			oldKey, err := crypto.CombinePasswordAndKey([]byte(masterpass), oldHex)
			check(err)

			key := crypto.GeneratePBKDF2Key([]byte(masterpass), []byte(privKey))

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()

			client := initClient(ctx)

			mdch, errch := client.GetAllMetadata(ctx, "")

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

					cred, err := client.Get(ctx, md.ID)
					check(err)

					cipherpass, err := base64.StdEncoding.DecodeString(cred.Password)
					check(err)

					plainpass, err := crypto.GCMDecrypt(oldKey, cipherpass)
					check(err)

					cipherpass, err = crypto.CBCEncrypt(key, plainpass)
					check(err)

					password := base64.StdEncoding.EncodeToString(cipherpass)

					var otpSecret string

					if cred.OTPSecret != "" {
						ciphersecret, err := base64.StdEncoding.DecodeString(cred.OTPSecret)
						check(err)

						plainsecret, err := crypto.GCMDecrypt(oldKey, ciphersecret)
						check(err)

						ciphersecret, err = crypto.CBCEncrypt(key, plainsecret)
						check(err)

						otpSecret = base64.StdEncoding.EncodeToString(ciphersecret)
					}

					credIn := types.CredentialInput{
						MetadataInput: types.MetadataInput{
							Primary:    cred.Primary,
							SourceHost: cred.SourceHost,
							LoginURL:   cred.LoginURL,
							Tag:        cred.Tag,
						},
						Username:  cred.Username,
						Email:     cred.Email,
						Password:  password,
						OTPSecret: otpSecret,
					}

					_, err = client.Update(ctx, cred.ID, credIn)
					check(err)
				}
			}
		},
	}

	return gcmToCBC
}
