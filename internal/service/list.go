package service

import (
	"bytes"
	"fmt"
	"github.com/yemingfeng/sdb/internal/store"
	"github.com/yemingfeng/sdb/internal/store/engine"
	"github.com/yemingfeng/sdb/internal/util"
	"math"
	"strconv"
	"strings"
)

const listKeyPrefixTemplate = "ll/%s"
const listKeyTemplate = listKeyPrefixTemplate + "/%d"
const listIdKeyPrefixTemplate = "li/%s/%s"
const listIdKeyTemplate = listIdKeyPrefixTemplate + "/%d"

func LRPush(key []byte, values [][]byte) (bool, error) {
	lock(LList, key)
	defer unlock(LList, key)

	batch := store.NewBatch()
	defer batch.Close()

	for _, value := range values {
		id := util.GetOrderingKey()
		if _, err := batch.Set(generateListKey(key, id), value); err != nil {
			return false, err
		}
		if _, err := batch.Set(generateListIdKey(key, value, id), value); err != nil {
			return false, err
		}
	}

	return batch.Commit()
}

func LLPush(key []byte, values [][]byte) (bool, error) {
	lock(LList, key)
	defer unlock(LList, key)

	batch := store.NewBatch()
	defer batch.Close()

	for i, value := range values {
		id := -util.GetOrderingKey() - int64(i)
		if _, err := batch.Set(generateListKey(key, id), value); err != nil {
			return false, err
		}
		if _, err := batch.Set(generateListIdKey(key, value, id), value); err != nil {
			return false, err
		}
	}

	return batch.Commit()
}

func LPop(key []byte, values [][]byte) (bool, error) {
	lock(LList, key)
	defer unlock(LList, key)

	batch := store.NewBatch()
	defer batch.Close()

	for i := range values {
		if err := store.Iterate(&engine.PrefixIteratorOption{Prefix: generateListIdPrefixKey(key, values[i])},
			func(storeKey []byte, storeValue []byte) error {
				if bytes.Equal(storeValue, values[i]) {
					if _, err := batch.Del(storeKey); err != nil {
						return err
					}

					infos := strings.Split(string(storeKey), "/")
					id, err := strconv.ParseInt(infos[len(infos)-1], 10, 64)
					if err != nil {
						return err
					}
					if _, err := batch.Del(generateListKey(key, id)); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
			return false, err
		}
	}

	return batch.Commit()
}

func LRange(key []byte, offset int32, limit uint32) ([][]byte, error) {
	index := int32(0)
	res := make([][]byte, limit)
	_ = store.Iterate(&engine.PrefixIteratorOption{
		Prefix: generateListPrefixKey(key), Offset: offset, Limit: limit},
		func(key []byte, value []byte) error {
			res[index] = value
			index++
			return nil
		})
	return res[0:index], nil
}

func LExist(key []byte, values [][]byte) ([]bool, error) {
	res := make([]bool, len(values))
	existMap := make(map[string]bool)
	for i := range values {
		_ = store.Iterate(&engine.PrefixIteratorOption{Prefix: generateListIdPrefixKey(key, values[i])},
			func(key []byte, value []byte) error {
				if bytes.Equal(value, values[i]) {
					existMap[string(value)] = true
				}
				return nil
			})
	}
	for i, value := range values {
		if existMap[string(value)] {
			res[i] = true
		}
	}
	return res, nil
}

func LDel(key []byte) (bool, error) {
	lock(LList, key)
	defer unlock(LList, key)

	batch := store.NewBatch()
	defer batch.Close()

	if err := store.Iterate(&engine.PrefixIteratorOption{Prefix: generateListPrefixKey(key)},
		func(key1 []byte, value1 []byte) error {
			if _, err := batch.Del(key1); err != nil {
				return err
			}

			return store.Iterate(&engine.PrefixIteratorOption{Prefix: generateListIdPrefixKey(key, value1)},
				func(key2 []byte, value2 []byte) error {
					if bytes.Equal(value2, value1) {
						_, err := batch.Del(key2)
						return err
					}
					return nil
				})
		}); err != nil {
		return false, err
	}

	return batch.Commit()
}

func LCount(key []byte) (uint32, error) {
	count := uint32(0)
	_ = store.Iterate(&engine.PrefixIteratorOption{Prefix: generateListPrefixKey(key)},
		func(key []byte, value []byte) error {
			count++
			return nil
		})
	return count, nil
}

func LMembers(key []byte) ([][]byte, error) {
	index := int32(0)
	res := make([][]byte, 0)
	_ = store.Iterate(&engine.PrefixIteratorOption{
		Prefix: generateListPrefixKey(key), Offset: 0, Limit: math.MaxInt32},
		func(key []byte, value []byte) error {
			res = append(res, value)
			index++
			return nil
		})
	return res[0:index], nil
}

func generateListKey(key []byte, id int64) []byte {
	return []byte(fmt.Sprintf(listKeyTemplate, key, id))
}

func generateListPrefixKey(key []byte) []byte {
	return []byte(fmt.Sprintf(listKeyPrefixTemplate, key))
}

func generateListIdKey(key []byte, value []byte, id int64) []byte {
	return []byte(fmt.Sprintf(listIdKeyTemplate, key, value, id))
}

func generateListIdPrefixKey(key []byte, value []byte) []byte {
	return []byte(fmt.Sprintf(listIdKeyPrefixTemplate, key, value))
}
