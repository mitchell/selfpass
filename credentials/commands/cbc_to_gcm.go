package commands

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	clitypes "github.com/mitchell/selfpass/cli/types"
	"github.com/mitchell/selfpass/credentials/types"
	"github.com/mitchell/selfpass/crypto"
)

func MakeCBCtoGCM(repo clitypes.ConfigRepo, initClient CredentialClientInit) *cobra.Command {
	cbcToGCM := &cobra.Command{
		Use:    "cbc-to-gcm",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			masterpass, cfg, err := repo.OpenConfig()
			check(err)

			key, err := hex.DecodeString(cfg.GetString(clitypes.KeyPrivateKey))
			check(err)

			keypass, err := crypto.CombinePasswordAndKey([]byte(masterpass), key)
			check(err)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
			defer cancel()

			client := initClient(ctx)

			mdch, errch := client.GetAllMetadata(ctx, "")

			for {
				select {
				case err := <-errch:
					check(err)
				case md, ok := <-mdch:
					if !ok {
						fmt.Println("All done.")
						return
					}

					cred, err := client.Get(ctx, md.ID)
					check(err)

					passbytes, err := base64.StdEncoding.DecodeString(cred.Password)
					check(err)

					plainpass, err := crypto.CBCDecrypt(keypass, passbytes)
					check(err)

					passbytes, err = crypto.GCMEncrypt(keypass, plainpass)
					check(err)

					cred.Password = base64.StdEncoding.EncodeToString(passbytes)

					if cred.OTPSecret != "" {
						passbytes, err := base64.StdEncoding.DecodeString(cred.OTPSecret)
						check(err)

						plainpass, err := crypto.CBCDecrypt(keypass, passbytes)
						check(err)

						passbytes, err = crypto.GCMEncrypt(keypass, plainpass)
						check(err)

						cred.OTPSecret = base64.StdEncoding.EncodeToString(passbytes)
					}

					_, err = client.Update(ctx, cred.ID, types.CredentialInput{
						MetadataInput: types.MetadataInput{
							Tag:        cred.Tag,
							SourceHost: cred.SourceHost,
							LoginURL:   cred.LoginURL,
							Primary:    cred.Primary,
						},
						OTPSecret: cred.OTPSecret,
						Password:  cred.Password,
						Email:     cred.Email,
						Username:  cred.Username,
					})
					check(err)
				}
			}
		},
	}

	return cbcToGCM
}
