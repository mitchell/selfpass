package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/mitchell/selfpass/services/cli/types"
	"github.com/mitchell/selfpass/services/crypto"
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

			key := cfg.GetString(types.KeyPrivateKey)
			passkey := crypto.GeneratePBKDF2Key([]byte(masterpass), []byte(key))

			contents, err = crypto.GCMEncrypt(passkey, contents)
			check(err)

			check(ioutil.WriteFile(fileEnc, contents, 0600))
			check(os.Remove(file))

			fmt.Println("Encrypted file: ", fileEnc)
		},
	}

	return encryptCmd
}
