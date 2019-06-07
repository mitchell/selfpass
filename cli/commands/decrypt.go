package commands

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/mitchell/selfpass/cli/types"
	"github.com/mitchell/selfpass/crypto"
)

func makeDecrypt(repo types.ConfigRepo) *cobra.Command {
	decryptCmd := &cobra.Command{
		Use:   "decrypt [file]",
		Short: "Decrypt a file using your masterpass and secret key",
		Long: `Decrypt a file using your masterpass and secret key, and replace the old file with
the new file.`,
		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			masterpass, cfg, err := repo.OpenConfig()
			check(err)

			file := args[0]
			fileout := file

			if file[len(file)-4:] == ".enc" {
				fileout = file[:len(file)-4]
			}

			contents, err := ioutil.ReadFile(file)
			check(err)

			key, err := hex.DecodeString(cfg.GetString(types.KeyPrivateKey))
			check(err)

			passkey, err := crypto.CombinePasswordAndKey([]byte(masterpass), []byte(key))
			check(err)

			contents, err = crypto.GCMDecrypt(passkey, contents)
			check(err)

			check(ioutil.WriteFile(fileout, contents, 0600))
			check(os.Remove(file))

			fmt.Println("Decrypted file: ", fileout)
		},
	}

	return decryptCmd
}
