package service

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/store"
)

const setKeyPrefixTemplate = "sk/%s"
const setKeyTemplate = setKeyPrefixTemplate + "/%s"

func SPush(key []byte, values [][]byte) (bool, error) {
	lock(LSet, key)
	defer unlock(LSet, key)

	batch := store.NewBatch()
	defer batch.Close()

	for _, value := range values {
		if _, err := batch.Set(generateSetKey(key, value), value); err != nil {
			return false, err
		}
	}

	return batch.Commit()
}

func SPop(key []byte, values [][]byte) (bool, error) {
	lock(LSet, key)
	defer unlock(LSet, key)

	batch := store.NewBatch()
	defer batch.Close()

	for _, value := range values {
		if _, err := batch.Del(generateSetKey(key, value)); err != nil {
			return false, err
		}
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

	if err := store.Iterate0(generateSetPrefixKey(key),
		func(key []byte, value []byte) error {
			_, err := batch.Del(key)
			return err
		}); err != nil {
		return false, err
	}

	return batch.Commit()
}

func SCount(key []byte) (uint32, error) {
	count := uint32(0)
	_ = store.Iterate0(generateSetPrefixKey(key),
		func(key []byte, value []byte) error {
			count = count + 1
			return nil
		})
	return count, nil
}

func SMembers(key []byte) ([][]byte, error) {
	index := int32(0)
	res := make([][]byte, 0)
	_ = store.Iterate0(generateSetPrefixKey(key),
		func(key []byte, value []byte) error {
			res = append(res, value)
			index++
			return nil
		})
	return res[0:index], nil
}

func generateSetPrefixKey(key []byte) []byte {
	return []byte(fmt.Sprintf(setKeyPrefixTemplate, key))
}

func generateSetKey(key []byte, value []byte) []byte {
	return []byte(fmt.Sprintf(setKeyTemplate, key, value))
}
