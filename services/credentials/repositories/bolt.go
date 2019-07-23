package repositories

import (
	"bytes"
	"context"
	"encoding/gob"
	"os"
	"sync"

	"go.etcd.io/bbolt"

	"github.com/mitchell/selfpass/services/credentials/types"
)

func OpenBoltDB(file string, mode os.FileMode, opts *bbolt.Options) (out BoltDB, err error) {
	db, err := bbolt.Open(file, mode, opts)
	if err != nil {
		return out, err
	}

	return BoltDB{bolt: db}, nil
}

type BoltDB struct {
	bolt *bbolt.DB
}

func (db BoltDB) GetAllMetadata(ctx context.Context, sourceHost string, errch chan<- error) (output <-chan types.Metadata) {
	mdch := make(chan types.Metadata, 1)

	go func() {
		defer close(mdch)

		err := db.bolt.View(func(tx *bbolt.Tx) error {
			bkt := tx.Bucket([]byte(credentialsBkt))
			if bkt == nil {
				return nil
			}

			var wg sync.WaitGroup
			err := bkt.ForEach(func(_, value []byte) error {
				wg.Add(1)

				go func(value []byte) {
					defer wg.Done()

					var cred types.Credential

					err := gobUnmarshal(value, &cred)
					if err != nil {
						errch <- err
						return
					}

					if sourceHost == "" || sourceHost == cred.SourceHost {
						mdch <- cred.Metadata
					}
				}(value)

				return nil
			})
			if err != nil {
				return err
			}

			wg.Wait()

			return nil
		})
		if err != nil {
			errch <- err
			return
		}
	}()

	return mdch
}

func (db BoltDB) Get(ctx context.Context, id string) (output types.Credential, err error) {
	err = db.bolt.View(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(credentialsBkt))
		if bkt == nil {
			return nil
		}

		value := bkt.Get([]byte(id))
		if value == nil {
			return nil
		}

		return gobUnmarshal(value, &output)
	})

	return output, err
}

func (db BoltDB) Put(ctx context.Context, c types.Credential) (err error) {
	err = db.bolt.Update(func(tx *bbolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(credentialsBkt))
		if err != nil {
			return err
		}

		value, err := gobMarshal(c)
		if err != nil {
			return err
		}

		return bkt.Put([]byte(c.ID), value)
	})

	return err
}

func (db BoltDB) Delete(ctx context.Context, id string) (err error) {
	err = db.bolt.Update(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(credentialsBkt))
		if bkt == nil {
			return nil
		}

		return bkt.Delete([]byte(id))
	})

	return err
}

const credentialsBkt = "credentials"

func gobMarshal(v interface{}) (bs []byte, err error) {
	buf := bytes.NewBuffer(nil)
	err = gob.NewEncoder(buf).Encode(v)
	return buf.Bytes(), err
}

func gobUnmarshal(bs []byte, v interface{}) error {
	buf := bytes.NewReader(bs)
	return gob.NewDecoder(buf).Decode(v)
}
