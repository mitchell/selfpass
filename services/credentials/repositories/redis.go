package repositories

import (
	"context"

	"github.com/mediocregopher/radix/v3"
	"github.com/mitchell/selfpass/services/credentials/types"
)

func NewRedisConn(networkType, address string, connCount uint, options ...radix.PoolOpt) (c RedisConn, err error) {
	p, err := radix.NewPool(networkType, address, int(connCount), options...)
	return RedisConn{p: p}, err
}

type RedisConn struct {
	p *radix.Pool
}

func (conn RedisConn) GetAllMetadata(ctx context.Context, sourceHost string, errch chan<- error) (output <-chan types.Metadata) {
	mdch := make(chan types.Metadata, 1)

	go func() {
		defer close(mdch)

		var key string
		scr := radix.NewScanner(conn.p, radix.ScanOpts{Command: scan, Pattern: types.TypePrefixCred + dash + sourceHost + star})

		for scr.Next(&key) {
			select {
			case <-ctx.Done():
				return
			default:
			}

			var md types.Metadata

			if err := conn.p.Do(radix.Cmd(&md, hGetAll, key)); err != nil {
				errch <- err
				return
			}

			mdch <- md
		}
	}()

	return mdch
}

func (conn RedisConn) Get(ctx context.Context, id string) (output types.Credential, err error) {
	err = conn.p.Do(radix.Cmd(&output, hGetAll, id))
	return output, err
}

func (conn RedisConn) Put(ctx context.Context, c types.Credential) (err error) {
	err = conn.p.Do(radix.FlatCmd(nil, hMSet, c.ID, c))
	return err
}

func (conn RedisConn) Delete(ctx context.Context, id string) (err error) {
	err = conn.p.Do(radix.Cmd(nil, del, id))
	return err
}

func (conn RedisConn) DumpDB(ctx context.Context) (bs []byte, err error) {
	bs = []byte{}

	if err := conn.p.Do(radix.Cmd(&bs, "DUMP")); err != nil {
		return nil, err
	}

	return bs, nil
}

const (
	dash    = "-"
	star    = "*"
	scan    = "SCAN"
	hGetAll = "HGETALL"
	hMSet   = "HMSET"
	del     = "DEL"
)
