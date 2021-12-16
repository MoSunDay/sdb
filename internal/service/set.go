package service

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/store"
	"github.com/yemingfeng/sdb/internal/store/engine"
	"math"
)

const setKeyPrefixTemplate = "sk/%s"
const setKeyTemplate = setKeyPrefixTemplate + "/%s"

func SPush(key []byte, values [][]byte) (bool, error) {
	lock(LSet, key)
	defer unlock(LSet, key)

	batch := store.NewBatch()
	defer batch.Close()

	for _, value := range values {
		batch.Set(generateSetKey(key, value), value)
	}

	return batch.Commit()
}

func SPop(key []byte, values [][]byte) (bool, error) {
	lock(LSet, key)
	defer unlock(LSet, key)

	batch := store.NewBatch()
	defer batch.Close()

	for _, value := range values {
		batch.Del(generateSetKey(key, value))
	}

	return batch.Commit()
}

func SExist(key []byte, values [][]byte) ([]bool, error) {
	res := make([]bool, len(values))
	for i, value := range values {
		exist, err := store.Get(generateSetKey(key, value))
		if err != nil {
			return nil, err
		}
		res[i] = len(exist) > 0
	}
	return res, nil
}

func SDel(key []byte) (bool, error) {
	lock(LSet, key)
	defer unlock(LSet, key)

	batch := store.NewBatch()
	defer batch.Close()

	store.Iterate(&engine.PrefixIteratorOption{Prefix: generateSetPrefixKey(key)},
		func(key []byte, value []byte) {
			batch.Del(key)
		})

	return batch.Commit()
}

func SCount(key []byte) (uint32, error) {
	count := uint32(0)
	store.Iterate(&engine.PrefixIteratorOption{Prefix: generateSetPrefixKey(key)},
		func(key []byte, value []byte) {
			count = count + 1
		})
	return count, nil
}

func SMembers(key []byte) ([][]byte, error) {
	index := int32(0)
	res := make([][]byte, 0)
	store.Iterate(&engine.PrefixIteratorOption{
		Prefix: generateSetPrefixKey(key), Offset: 0, Limit: math.MaxInt32},
		func(key []byte, value []byte) {
			res = append(res, value)
			index++
		})
	return res[0:index], nil
}

func generateSetPrefixKey(key []byte) []byte {
	return []byte(fmt.Sprintf(setKeyPrefixTemplate, key))
}

func generateSetKey(key []byte, value []byte) []byte {
	return []byte(fmt.Sprintf(setKeyTemplate, key, value))
}
