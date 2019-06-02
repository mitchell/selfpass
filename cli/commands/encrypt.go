package commands

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/mitchell/selfpass/cli/types"
	"github.com/mitchell/selfpass/credentials/commands"
	"github.com/mitchell/selfpass/crypto"
)

func makeEncrypt(repo types.ConfigRepo) *cobra.Command {
	encryptCmd := &cobra.Command{
		Use:   "encrypt [file]",
		Short: "Encrypt a file using your masterpass and secret key",
		Long: `Encrypt a file using your masterpass and secret key, and replace the old file with the
new file.`,
		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			masterpass, cfg, err := repo.OpenConfig()
			check(err)

			file := args[0]
			fileEnc := file + ".enc"

			contents, err := ioutil.ReadFile(file)
			check(err)

			key, err := hex.DecodeString(cfg.GetString(commands.KeyPrivateKey))
			check(err)

			passkey, err := crypto.CombinePasswordAndKey([]byte(masterpass), []byte(key))
			check(err)

			contents, err = crypto.CBCEncrypt(passkey, contents)
			check(err)

			check(ioutil.WriteFile(fileEnc, contents, 0600))
			check(os.Remove(file))

			fmt.Println("Encrypted file: ", fileEnc)
		},
	}

	return encryptCmd
}
