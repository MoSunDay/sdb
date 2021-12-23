package badger

import (
	"github.com/dgraph-io/badger/v3"
)

type BadgerBatch struct {
	transaction *badger.Txn
}

func (batch *BadgerBatch) Get(key []byte) ([]byte, error) {
	item, err := batch.transaction.Get(key)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return item.ValueCopy(nil)
}

func (batch *BadgerBatch) Set(key []byte, value []byte) (bool, error) {
	if err := batch.transaction.Set(key, value); err != nil {
		return false, err
	}
	return true, nil
}

func (batch *BadgerBatch) Del(key []byte) (bool, error) {
	if err := batch.transaction.Delete(key); err != nil {
		return false, err
	}
	return true, nil
}

func (batch *BadgerBatch) Commit() (bool, error) {
	if err := batch.transaction.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func (batch *BadgerBatch) Close() {
	batch.transaction.Discard()
}
