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

const listKeyPrefixTemplate = "l/%s"
const listKeyTemplate = listKeyPrefixTemplate + "/%d"
const listIdKeyPrefixTemplate = "li/%s/%s"
const listIdKeyTemplate = listIdKeyPrefixTemplate + "/%d"

func LPush(key []byte, values [][]byte) (bool, error) {
	lock(LList, key)
	defer unlock(LList, key)

	batch := store.NewBatch()
	defer batch.Close()

	for _, value := range values {
		id := util.GetOrderingKey()
		batch.Set(generateListKey(key, id), value)
		batch.Set(generateListIdKey(key, value, id), value)
	}

	return batch.Commit()
}

func LPop(key []byte, values [][]byte) (bool, error) {
	lock(LList, key)
	defer unlock(LList, key)

	batch := store.NewBatch()
	defer batch.Close()

	for i := range values {
		store.Iterate(&engine.PrefixIteratorOption{Prefix: generateListIdPrefixKey(key, values[i])},
			func(storeKey []byte, storeValue []byte) {
				if bytes.Equal(storeValue, values[i]) {
					batch.Del(storeKey)

					infos := strings.Split(string(storeKey), "/")
					id, _ := strconv.ParseInt(infos[len(infos)-1], 10, 64)
					batch.Del(generateListKey(key, id))
				}
			})
	}

	return batch.Commit()
}

func LRange(key []byte, offset int32, limit int32) ([][]byte, error) {
	index := int32(0)
	res := make([][]byte, limit)
	store.Iterate(&engine.PrefixIteratorOption{
		Prefix: generateListPrefixKey(key), Offset: int(offset), Limit: int(limit)},
		func(key []byte, value []byte) {
			res[index] = value
			index++
		})
	return res[0:index], nil
}

func LExist(key []byte, values [][]byte) ([]bool, error) {
	res := make([]bool, len(values))
	existMap := make(map[string]bool)
	for i := range values {
		store.Iterate(&engine.PrefixIteratorOption{Prefix: generateListIdPrefixKey(key, values[i])},
			func(key []byte, value []byte) {
				if bytes.Equal(value, values[i]) {
					existMap[string(value)] = true
				}
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

	store.Iterate(&engine.PrefixIteratorOption{Prefix: generateListPrefixKey(key)},
		func(key1 []byte, value1 []byte) {
			batch.Del(key1)

			store.Iterate(&engine.PrefixIteratorOption{Prefix: generateListIdPrefixKey(key, value1)},
				func(key2 []byte, value2 []byte) {
					if bytes.Equal(value2, value1) {
						batch.Del(key2)
					}
				})
		})

	return batch.Commit()
}

func LCount(key []byte) (int32, error) {
	count := int32(0)
	store.Iterate(&engine.PrefixIteratorOption{Prefix: generateListPrefixKey(key)},
		func(key []byte, value []byte) {
			count++
		})
	return count, nil
}

func LMembers(key []byte) ([][]byte, error) {
	index := int32(0)
	res := make([][]byte, 0)
	store.Iterate(&engine.PrefixIteratorOption{
		Prefix: generateListPrefixKey(key), Offset: 0, Limit: math.MaxInt32},
		func(key []byte, value []byte) {
			res = append(res, value)
			index++
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
