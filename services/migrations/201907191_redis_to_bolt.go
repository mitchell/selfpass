package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/mediocregopher/radix/v3"
	"github.com/spf13/pflag"
	"go.etcd.io/bbolt"

	"github.com/mitchell/selfpass/services/credentials/types"
	"github.com/mitchell/selfpass/services/migrations/migration"
)

const keyCredentials = "credentials"

func main() {
	redisHost := pflag.StringP("redis-host", "r", "127.0.0.1:6379", "specify the redis host")
	boltFile := pflag.StringP("bolt-file", "b", "./data/bolt.db", "specify the bolt DB file")
	help := pflag.BoolP("help", "h", false, "see help")
	pflag.Parse()

	if *help {
		pflag.PrintDefaults()
		return
	}

	pool, err := radix.NewPool("tcp", *redisHost, 10)
	migration.Check(err)

	db, err := bbolt.Open(*boltFile, 0600, nil)
	migration.Check(err)

	defer func() { migration.Check(db.Close()); migration.Check(pool.Close()) }()

	fmt.Println("Beginning migration...")

	var wg sync.WaitGroup
	scanner := radix.NewScanner(pool, radix.ScanOpts{Command: "SCAN"})

	for key := ""; scanner.Next(&key); {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()

			var cred types.Credential
			migration.Check(pool.Do(radix.Cmd(&cred, "HGETALL", key)))

			fmt.Printf("Migrating %s.\n", cred.ID)

			migration.Check(db.Batch(func(tx *bbolt.Tx) error {
				credBkt, err := tx.CreateBucketIfNotExists([]byte(keyCredentials))
				if err != nil {
					return err
				}

				buf := bytes.NewBuffer(nil)
				err = gob.NewEncoder(buf).Encode(cred)
				if err != nil {
					return err
				}

				if err = credBkt.Put([]byte(cred.ID), buf.Bytes()); err != nil {
					return err
				}

				return nil
			}))
		}(key)
	}

	wg.Wait()

	fmt.Println("Done migrating.")
}
