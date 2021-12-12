package badger

import (
	"github.com/dgraph-io/badger/v3"
)

type BadgerBatch struct {
	batch *badger.WriteBatch
}

func (batch *BadgerBatch) Set(key []byte, value []byte) {
	_ = batch.batch.Set(key, value)
}

func (batch *BadgerBatch) Del(key []byte) {
	_ = batch.batch.Delete(key)
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
