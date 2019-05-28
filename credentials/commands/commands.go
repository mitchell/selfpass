package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/mitchell/selfpass/credentials/types"
)

type CredentialClientInit func(ctx context.Context) (c types.CredentialClient)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const KeyPrivateKey = "private_key"
