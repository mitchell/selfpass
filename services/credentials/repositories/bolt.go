package repositories

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
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
			bkt := getCredentialsBucket(tx)
			if bkt.isEmpty {
				return nil
			}

			var wg sync.WaitGroup
			c := bkt.hostPrimaryIndex.Cursor()

			if sourceHost == "" {
				for key, value := c.First(); key != nil; key, value = c.Next() {
					wg.Add(1)
					unmarshalAndSendCred(value, mdch, errch, &wg)
				}
			} else {
				hostBytes := []byte(sourceHost)
				for key, value := c.Seek(hostBytes); bytes.HasPrefix(key, hostBytes); key, value = c.Next() {
					wg.Add(1)
					unmarshalAndSendCred(value, mdch, errch, &wg)
				}
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

func unmarshalAndSendCred(value []byte, mdch chan<- types.Metadata, errch chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	var cred types.Credential

	err := gobUnmarshal(value, &cred)
	if err != nil {
		errch <- err
		return
	}

	mdch <- cred.Metadata
}

func (db BoltDB) Get(ctx context.Context, id string) (output types.Credential, err error) {
	err = db.bolt.View(func(tx *bbolt.Tx) error {
		bkt := getCredentialsBucket(tx)
		if bkt.isEmpty {
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
	return db.bolt.Update(func(tx *bbolt.Tx) error {
		bkt := getCredentialsBucket(tx)
		bkt.createIfNotExists()

		value := bkt.Get([]byte(c.ID))
		if value != nil {
			var cred types.Credential
			if err = gobUnmarshal(value, &cred); err != nil {
				return err
			}

			if err = bkt.Delete([]byte(c.ID)); err != nil {
				return err
			}
			if err = bkt.hostPrimaryIndex.Delete([]byte(genHostPrimaryIdxKey(cred))); err != nil {
				return err
			}
		}

		value, err := gobMarshal(c)
		if err != nil {
			return err
		}

		if err = bkt.hostPrimaryIndex.Put([]byte(genHostPrimaryIdxKey(c)), value); err != nil {
			return err
		}

		return bkt.Put([]byte(c.ID), value)
	})
}

func (db BoltDB) Delete(ctx context.Context, id string) (err error) {
	return db.bolt.Update(func(tx *bbolt.Tx) error {
		bkt := getCredentialsBucket(tx)
		if bkt.isEmpty {
			return nil
		}

		value := bkt.Get([]byte(id))
		if value == nil {
			return nil
		}

		var cred types.Credential
		if err = gobUnmarshal(value, &cred); err != nil {
			return err
		}

		if err = bkt.hostPrimaryIndex.Delete([]byte(genHostPrimaryIdxKey(cred))); err != nil {
			return err
		}

		return bkt.Delete([]byte(id))
	})
}

const keyCredentialsBkt = "credentials"
const keyHostAndPrimaryIdx = "sourceHost-primary"

func getCredentialsBucket(tx *bbolt.Tx) credentialsBucket {
	bkt := credentialsBucket{
		Bucket: tx.Bucket([]byte(keyCredentialsBkt)),
		tx:     tx,
	}
	bkt.isEmpty = bkt.Bucket == nil

	if !bkt.isEmpty {
		bkt.hostPrimaryIndex = bkt.Bucket.Bucket([]byte(keyHostAndPrimaryIdx))
	}

	return bkt
}

type credentialsBucket struct {
	*bbolt.Bucket
	tx               *bbolt.Tx
	hostPrimaryIndex *bbolt.Bucket
	isEmpty          bool
}

func (bkt *credentialsBucket) createIfNotExists() {
	if bkt.isEmpty {
		bkt.Bucket, _ = bkt.tx.CreateBucket([]byte(keyCredentialsBkt))
		bkt.hostPrimaryIndex, _ = bkt.CreateBucket([]byte(keyHostAndPrimaryIdx))
		bkt.isEmpty = false
	}
}

func genHostPrimaryIdxKey(cred types.Credential) string {
	return fmt.Sprintf("%s-%s-%s", cred.SourceHost, cred.Primary, cred.ID)
}

func gobMarshal(v interface{}) (bs []byte, err error) {
	buf := bytes.NewBuffer(nil)
	err = gob.NewEncoder(buf).Encode(v)
	return buf.Bytes(), err
}

func gobUnmarshal(bs []byte, v interface{}) error {
	buf := bytes.NewReader(bs)
	return gob.NewDecoder(buf).Decode(v)
}
