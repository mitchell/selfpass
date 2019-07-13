# sp

This is the project home of `sp` the *selfpass* CLI.

To install with Go run `go get -u github.com/mitchell/selfpass/sp`.

Help menu:
```
This is the CLI client for Selfpass, the self-hosted password manager. With this tool you
can interact with the entire Selfpass API.

Usage:
  sp [command]

Available Commands:
  create      Create a credential in Selfpass
  decrypt     Decrypt a file using your masterpass and secret key
  decrypt-cfg Decrypt your config file
  delete      Delete a credential using the given ID
  encrypt     Encrypt a file using your masterpass and secret key
  get         Get a credential info and copy password to clipboard
  help        Help about any command
  init        This command initializes SPC for the first time
  list        List the metadata for all credentials
  update      Update a credential in Selfpass

Flags:
      --config string   config file (default is $HOME/.sp.toml)
  -h, --help            help for sp
      --version         version for sp

Use "sp [command] --help" for more information about a command.
```

For more project-level information see the root `README.md`.