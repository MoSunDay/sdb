package pebble

import "github.com/cockroachdb/pebble"

type PebbleBatch struct {
	batch *pebble.Batch
}

func (batch *PebbleBatch) Set(key []byte, value []byte) (bool, error) {
	if err := batch.batch.Set(key, value, nil); err != nil {
		return false, err
	}
	return true, nil
}

func (batch *PebbleBatch) Del(key []byte) (bool, error) {
	if err := batch.batch.Delete(key, nil); err != nil {
		return false, err
	}
	return true, nil
}

func (batch *PebbleBatch) Commit() (bool, error) {
	if err := batch.batch.Commit(pebble.Sync); err != nil {
		return false, err
	}
	return true, nil
}

func (batch *PebbleBatch) Close() {
	_ = batch.batch.Close()
}
