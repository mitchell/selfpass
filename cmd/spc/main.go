package main

import (
	"context"

	"github.com/mitchell/selfpass/cmd/spc/cmd"
	"github.com/mitchell/selfpass/credentials/repositories"
)

func main() {
	ctx := context.Background()
	cmd.Execute(ctx, repositories.NewCredentialServiceClient)
}
