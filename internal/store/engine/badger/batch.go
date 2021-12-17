package badger

import (
	"github.com/dgraph-io/badger/v3"
)

type BadgerBatch struct {
	batch *badger.WriteBatch
}

func (batch *BadgerBatch) Set(key []byte, value []byte) (bool, error) {
	if err := batch.batch.Set(key, value); err != nil {
		return false, err
	}
	return true, nil
}

func (batch *BadgerBatch) Del(key []byte) (bool, error) {
	if err := batch.batch.Delete(key); err != nil {
		return false, err
	}
	return true, nil
}

func (batch *BadgerBatch) Commit() (bool, error) {
	if err := batch.batch.Flush(); err != nil {
		return false, err
	}
	return true, nil
}

func (batch *BadgerBatch) Close() {
	batch.batch.Cancel()
}
