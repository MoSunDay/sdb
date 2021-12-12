package level

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type LevelBatch struct {
	db    *leveldb.DB
	batch *leveldb.Batch
}

func (batch *LevelBatch) Set(key []byte, value []byte) {
	batch.batch.Put(key, value)
}

func (batch *LevelBatch) Del(key []byte) {
	batch.batch.Delete(key)
}

func (batch *LevelBatch) Commit() (bool, error) {
	if err := batch.db.Write(batch.batch, &opt.WriteOptions{Sync: true}); err != nil {
		return false, err
	}
	return true, nil
}

func (batch *LevelBatch) Close() {
}
