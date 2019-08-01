package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"sync"

	"github.com/spf13/pflag"
	"go.etcd.io/bbolt"

	"github.com/mitchell/selfpass/services/credentials/types"
	"github.com/mitchell/selfpass/services/migrations/migration"
)

const keyCredentialsBkt = "credentials"
const keyHostAndPrimaryIdx = "sourceHost-primary"

func main() {
	file := pflag.StringP("file", "f", "./data/bolt.db", "specify the bolt db file")
	help := pflag.BoolP("help", "h", false, "see help")
	pflag.Parse()

	if *help {
		pflag.PrintDefaults()
		return
	}

	db, err := bbolt.Open(*file, 0600, nil)
	migration.Check(err)

	fmt.Println("Beginning migration...")

	creds := make(chan types.Credential, 1)
	errs := make(chan error, 1)

	go func() {
		defer close(creds)

		var wg sync.WaitGroup

		errs <- db.View(func(tx *bbolt.Tx) error {
			bkt := tx.Bucket([]byte(keyCredentialsBkt))
			if bkt == nil {
				return errors.New("no credentials bucket")
			}

			return bkt.ForEach(func(_, value []byte) error {
				wg.Add(1)

				go func(value []byte) {
					defer wg.Done()

					reader := bytes.NewReader(value)

					var cred types.Credential
					errs <- gob.NewDecoder(reader).Decode(&cred)

					creds <- cred
				}(value)

				return nil
			})
		})

		wg.Wait()
	}()

	go func() {
		defer close(errs)

		var wg sync.WaitGroup

		for cred := range creds {
			key := fmt.Sprintf("%s-%s-%s", cred.SourceHost, cred.Primary, cred.ID)

			fmt.Printf("Adding credential %s to index as %s.\n", cred.ID, key)

			wg.Add(1)
			go func(key string, cred types.Credential) {
				defer wg.Done()

				buf := bytes.NewBuffer(nil)
				migration.Check(gob.NewEncoder(buf).Encode(cred))

				value := buf.Bytes()

				errs <- db.Batch(func(tx *bbolt.Tx) error {
					credBkt := tx.Bucket([]byte(keyCredentialsBkt))

					bkt, err := credBkt.CreateBucketIfNotExists([]byte(keyHostAndPrimaryIdx))
					if err != nil {
						return err
					}

					return bkt.Put([]byte(key), value)
				})
			}(key, cred)
		}

		wg.Wait()
	}()

	for err = range errs {
		migration.Check(err)
	}

	fmt.Println("Migration done.")
}
