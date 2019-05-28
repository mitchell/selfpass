package main

import (
	"github.com/mitchell/selfpass/cli/commands"
	"github.com/mitchell/selfpass/credentials/repositories"
)

func main() {
	commands.Execute(repositories.NewCredentialServiceClient)
}
