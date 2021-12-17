package main

import (
	"github.com/yemingfeng/sdb/internal/store"
	"github.com/yemingfeng/sdb/internal/store/engine"
	"log"
	"strconv"
)

func main() {
	for i := 0; i < 10; i++ {
		_, _ = store.Set([]byte("h"+strconv.Itoa(i)), []byte("w"+strconv.Itoa(i)))
	}
	_ = store.Iterate(&engine.PrefixIteratorOption{Prefix: []byte("h"), Offset: -1, Limit: 3},
		func(key []byte, value []byte) error {
			log.Printf("key = %s, value = %s", key, value)
			return nil
		})
	log.Printf("=====")
	_ = store.Iterate(&engine.PrefixIteratorOption{Prefix: []byte("h"), Offset: -9, Limit: 3},
		func(key []byte, value []byte) error {
			log.Printf("key = %s, value = %s", key, value)
			return nil
		})
	log.Printf("=====")
	_ = store.Iterate(&engine.PrefixIteratorOption{Prefix: []byte("h"), Offset: 0, Limit: 3},
		func(key []byte, value []byte) error {
			log.Printf("key = %s, value = %s", key, value)
			return nil
		})
	log.Printf("=====")
	_ = store.Iterate(&engine.PrefixIteratorOption{Prefix: []byte("h"), Offset: 3, Limit: 5},
		func(key []byte, value []byte) error {
			log.Printf("key = %s, value = %s", key, value)
			return nil
		})
}
