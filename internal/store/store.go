package store

import (
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/store/engine"
	"github.com/yemingfeng/sdb/internal/store/engine/badger"
	"github.com/yemingfeng/sdb/internal/store/engine/level"
	"github.com/yemingfeng/sdb/internal/store/engine/pebble"
	"log"
	"math"
)

const (
	PEBBLE string = "pebble"
	BADGER string = "badger"
	LEVEL  string = "level"
)

var store engine.Store

func init() {
	if conf.Conf.Store.Engine == PEBBLE {
		store = pebble.NewPebbleStore()
	} else if conf.Conf.Store.Engine == BADGER {
		store = badger.NewBadgerStore()
	} else if conf.Conf.Store.Engine == LEVEL {
		store = level.NewLevelStore()
	} else {
		log.Fatalf("not match store engine: %s", conf.Conf.Store.Engine)
	}
}

func Set(key []byte, value []byte) (bool, error) {
	return store.Set(key, value)
}

func Get(key []byte) ([]byte, error) {
	return store.Get(key)
}

func Del(key []byte) (bool, error) {
	return store.Del(key)
}

func NewBatch() engine.Batch {
	return store.NewBatch()
}

func Iterate0(prefix []byte, handle func([]byte, []byte) error) error {
	return Iterate1(prefix, 0, math.MaxUint32, handle)
}

func Iterate1(prefix []byte, offset int32, limit uint32, handle func([]byte, []byte) error) error {
	if limit == 0 {
		limit = math.MaxUint32
	}
	return store.Iterate(&engine.PrefixIteratorOption{Prefix: prefix, Offset: offset, Limit: limit}, handle)
}

// Close todo call
func Close() error {
	return store.Close()
}
