package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/mitchell/selfpass/services/cli/types"
	"github.com/mitchell/selfpass/services/crypto"
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

			key := cfg.GetString(types.KeyPrivateKey)
			passkey := crypto.GeneratePBKDF2Key([]byte(masterpass), []byte(key))

			contents, err = crypto.GCMDecrypt(passkey, contents)
			check(err)

			check(ioutil.WriteFile(fileout, contents, 0600))
			check(os.Remove(file))

			fmt.Println("Decrypted file: ", fileout)
		},
	}

	return decryptCmd
}
