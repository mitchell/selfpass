package commands

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mitchell/selfpass/credentials/commands"
	"github.com/mitchell/selfpass/crypto"
)

func makeDecrypt(masterpass string, cfg *viper.Viper) *cobra.Command {
	decryptCmd := &cobra.Command{
		Use:   "decrypt [file]",
		Short: "Decrypt a file using your masterpass and secret key",
		Long: `Decrypt a file using your masterpass and secret key, and replace the old file with
the new file.`,
		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			file := args[0]
			fileout := file

			if file[len(file)-4:] == ".enc" {
				fileout = file[:len(file)-4]
			}

			contents, err := ioutil.ReadFile(file)
			check(err)

			key, err := hex.DecodeString(cfg.GetString(commands.KeyPrivateKey))
			check(err)

			passkey, err := crypto.CombinePasswordAndKey([]byte(masterpass), []byte(key))
			check(err)

			contents, err = crypto.CBCDecrypt(passkey, contents)
			check(err)

			check(ioutil.WriteFile(fileout, contents, 0600))
			check(os.Remove(file))

			fmt.Println("Decrypted file: ", fileout)
		},
	}

	return decryptCmd
}
