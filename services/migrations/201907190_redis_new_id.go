package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mediocregopher/radix/v3"
	"github.com/spf13/pflag"

	"github.com/mitchell/selfpass/services/credentials/types"
	"github.com/mitchell/selfpass/services/migrations/migration"
)

func main() {
	redisHost := pflag.StringP("redis-host", "r", "127.0.0.1:6379", "specify the redis host to target")
	help := pflag.BoolP("help", "h", false, "see help")
	pflag.Parse()

	if *help {
		pflag.PrintDefaults()
		return
	}

	pool, err := radix.NewPool("tcp", *redisHost, 10)
	migration.Check(err)

	fmt.Println("Beginning migration...")

	var pipeCmds []radix.CmdAction
	var creds []*types.Credential
	scanner := radix.NewScanner(pool, radix.ScanAllKeys)

	for key := ""; scanner.Next(&key); {
		var cred types.Credential
		pipeCmds = append(pipeCmds, radix.Cmd(&cred, "HGETALL", key))
		pipeCmds = append(pipeCmds, radix.Cmd(nil, "DEL", key))
		creds = append(creds, &cred)
	}

	migration.Check(pool.Do(radix.Pipeline(pipeCmds...)))
	pipeCmds = nil

	for _, cred := range creds {
		tcred := *cred
		tcred.ID = generateID()

		fmt.Printf("Migrating %s to %s.\n", cred.ID, tcred.ID)

		pipeCmds = append(pipeCmds, radix.FlatCmd(nil, "HMSET", tcred.ID, tcred))
	}

	migration.Check(pool.Do(radix.Pipeline(pipeCmds...)))

	fmt.Println("Done migrating.")
}

func generateID() string {
	const idLen = 8
	const alphanumerics = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ0123456789"
	const alphaLen = len(alphanumerics)

	rand.Seed(time.Now().UnixNano())
	id := make([]byte, idLen)

	for index := range id {
		id[index] = alphanumerics[rand.Int63()%int64(alphaLen)]
	}

	return fmt.Sprintf("%s-%s", types.KeyCredential, string(id))
}
