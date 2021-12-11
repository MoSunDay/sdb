package pebble

import (
	"github.com/cockroachdb/pebble"
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/store/engine"
	"github.com/yemingfeng/sdb/internal/util"
	"log"
)

type PebbleStore struct {
	db *pebble.DB
}

func NewPebbleStore() *PebbleStore {
	dbPath := conf.Conf.Store.Path + "/pebble"
	db, err := pebble.Open(dbPath, &pebble.Options{})
	if err != nil {
		log.Fatalf("failed to open file: %+v", err)
	}
	log.Printf("db init %s complete", dbPath)

	return &PebbleStore{db: db}
}

func (store *PebbleStore) Set(key []byte, value []byte, sync bool) (bool, error) {
	err := store.db.Set(key, value, &pebble.WriteOptions{Sync: sync})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (store *PebbleStore) Get(key []byte) ([]byte, error) {
	value, closer, err := store.db.Get(key)
	if err == pebble.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err = closer.Close(); err != nil {
		return nil, err
	}
	return value, err
}

func (store *PebbleStore) Del(key []byte, sync bool) (bool, error) {
	if err := store.db.Delete(key, &pebble.WriteOptions{Sync: sync}); err != nil {
		return false, nil
	}
	return true, nil
}

func (store *PebbleStore) NewBatch() engine.Batch {
	return &PebbleBatch{batch: store.db.NewBatch()}
}

func (store *PebbleStore) Iterate(opt *engine.PrefixIteratorOption, handle func([]byte, []byte)) {
	keyUpperBound := func(b []byte) []byte {
		end := make([]byte, len(b))
		copy(end, b)
		for i := len(end) - 1; i >= 0; i-- {
			end[i] = end[i] + 1
			if end[i] != 0 {
				return end[:i+1]
			}
		}
		return nil
	}

	var it = store.db.NewIter(&pebble.IterOptions{
		LowerBound: opt.Prefix,
		UpperBound: keyUpperBound(opt.Prefix),
	})
	defer func() {
		_ = it.Close()
	}()

	if opt.Offset >= 0 {
		i := 0
		for it.First(); i < opt.Offset && it.Valid(); it.Next() {
			i++
		}

		i = 0
		for ; it.Valid(); it.Next() {
			handle(util.Copy2(it.Key()), util.Copy2(it.Value()))
			i++
			if opt.Limit > 0 && i == opt.Limit {
				break
			}
		}
	} else {
		i := 0
		for it.Last(); i < -opt.Offset-1 && it.Valid(); it.Prev() {
			i++
		}

		i = 0
		for ; it.Valid(); it.Prev() {
			handle(util.Copy2(it.Key()), util.Copy2(it.Value()))
			i++
			if opt.Limit > 0 && i == opt.Limit {
				break
			}
		}
	}
}

func (store *PebbleStore) Close() error {
	return store.db.Close()
}
