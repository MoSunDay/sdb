package level

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type LevelBatch struct {
	transaction *leveldb.Transaction
}

func (batch *LevelBatch) Get(key []byte) ([]byte, error) {
	value, err := batch.transaction.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return value, err
}

func (batch *LevelBatch) Set(key []byte, value []byte) (bool, error) {
	if err := batch.transaction.Put(key, value, &opt.WriteOptions{Sync: true}); err != nil {
		return false, err
	}
	return true, nil
}

func (batch *LevelBatch) Del(key []byte) (bool, error) {
	if err := batch.transaction.Delete(key, &opt.WriteOptions{Sync: true}); err != nil {
		return false, err
	}
	return true, nil
}

func (batch *LevelBatch) Commit() (bool, error) {
	if err := batch.transaction.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func (batch *LevelBatch) Close() {
}
