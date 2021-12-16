package service

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/store"
	"github.com/yemingfeng/sdb/internal/store/engine"
	"github.com/yemingfeng/sdb/pkg/pb"
	"math"
	"strings"
)

const mapKeyPrefixTemplate = "mk/%s"
const mapKeyTemplate = mapKeyPrefixTemplate + "/%s"

func MPush(key []byte, pairs []*pb.Pair) (bool, error) {
	lock(LMap, key)
	defer unlock(LMap, key)

	batch := store.NewBatch()
	defer batch.Close()

	for i := range pairs {
		batch.Set(generateMapKey(key, pairs[i].Key), pairs[i].Value)
	}

	return batch.Commit()
}

func MPop(key []byte, keys [][]byte) (bool, error) {
	lock(LMap, key)
	defer unlock(LMap, key)

	batch := store.NewBatch()
	defer batch.Close()

	for i := range keys {
		batch.Del(generateMapKey(key, keys[i]))
	}

	return batch.Commit()
}

func MExist(key []byte, keys [][]byte) ([]bool, error) {
	res := make([]bool, len(keys))
	for i := range keys {
		exist, err := store.Get(generateMapKey(key, keys[i]))
		if err != nil {
			return nil, err
		}
		res[i] = len(exist) > 0
	}
	return res, nil
}

func MDel(key []byte) (bool, error) {
	lock(LMap, key)
	defer unlock(LMap, key)

	batch := store.NewBatch()
	defer batch.Close()

	store.Iterate(&engine.PrefixIteratorOption{Prefix: generateMapPrefixKey(key)},
		func(key []byte, value []byte) {
			batch.Del(key)
		})

	return batch.Commit()
}

func MCount(key []byte) (uint32, error) {
	count := uint32(0)
	store.Iterate(&engine.PrefixIteratorOption{Prefix: generateMapPrefixKey(key)},
		func(key []byte, value []byte) {
			count = count + 1
		})
	return count, nil
}

func MMembers(key []byte) ([]*pb.Pair, error) {
	index := int32(0)
	res := make([]*pb.Pair, 0)
	store.Iterate(&engine.PrefixIteratorOption{
		Prefix: generateMapPrefixKey(key), Offset: 0, Limit: math.MaxInt32},
		func(key []byte, value []byte) {
			infos := strings.Split(string(key), "/")
			res = append(res, &pb.Pair{Key: []byte(infos[2]), Value: value})
			index++
		})
	return res[0:index], nil
}

func generateMapPrefixKey(key []byte) []byte {
	return []byte(fmt.Sprintf(mapKeyPrefixTemplate, key))
}

func generateMapKey(key []byte, value []byte) []byte {
	return []byte(fmt.Sprintf(mapKeyTemplate, key, value))
}
