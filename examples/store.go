package main

import (
	"github.com/yemingfeng/sdb/internal/store"
	"log"
	"strconv"
)

func main() {
	for i := 0; i < 10; i++ {
		_, _ = store.Set([]byte("h"+strconv.Itoa(i)), []byte("w"+strconv.Itoa(i)))
	}
	_ = store.Iterate1([]byte("h"), -1, 3,
		func(key []byte, value []byte) error {
			log.Printf("key = %s, value = %s", key, value)
			return nil
		})
	log.Printf("=====")
	_ = store.Iterate1([]byte("h"), -9, 3,
		func(key []byte, value []byte) error {
			log.Printf("key = %s, value = %s", key, value)
			return nil
		})
	log.Printf("=====")
	_ = store.Iterate1([]byte("h"), 0, 3,
		func(key []byte, value []byte) error {
			log.Printf("key = %s, value = %s", key, value)
			return nil
		})
	log.Printf("=====")
	_ = store.Iterate1([]byte("h"), 3, 5,
		func(key []byte, value []byte) error {
			log.Printf("key = %s, value = %s", key, value)
			return nil
		})
}
