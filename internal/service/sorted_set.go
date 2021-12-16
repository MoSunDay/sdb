package service

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/store"
	"github.com/yemingfeng/sdb/internal/store/engine"
	"github.com/yemingfeng/sdb/pkg/pb"
	"google.golang.org/protobuf/proto"
	"math"
)

const sortedSetScoreKeyPrefixTemplate = "zs/%s"
const sortedSetScoreKeyTemplate = sortedSetScoreKeyPrefixTemplate + "/%s"
const sortedSetTupleKeyPrefixTemplate = "zt/%s"
const sortedSetTupleKeyTemplate = sortedSetTupleKeyPrefixTemplate + "/%s/%s"

func ZPush(key []byte, tuples []*pb.Tuple) (bool, error) {
	lock(LSortedSet, key)
	defer unlock(LSortedSet, key)

	// tuples -> [ {value: a, score: 1.0}, {value:b, score:1.1}, {value: c, score: 0.9} ]
	batch := store.NewBatch()
	defer batch.Close()

	for _, tuple := range tuples {
		value, err := proto.Marshal(tuple)
		if err != nil {
			return false, err
		}
		// get key z/{key}/{value} 获取 tuple
		exist, err := store.Get(generateSortedSetScoreKey(key, tuple.Value))
		if err != nil {
			return false, err
		}

		// remove key zs/{key}/{score}/{value}
		if exist != nil && len(exist) > 0 {
			existTuple := pb.Tuple{}
			err := proto.Unmarshal(exist, &existTuple)
			if err != nil {
				return false, err
			}
			batch.Del(generateSortedSetTupleKey(key, existTuple.Score, tuple.Value))
		}

		// add key z/{key}/{value} -> tuple
		batch.Set(generateSortedSetScoreKey(key, tuple.Value), value)

		// add key zs/{key}/{score}/{value} -> tuple
		batch.Set(generateSortedSetTupleKey(key, tuple.Score, tuple.Value), value)
	}

	return batch.Commit()
}

func ZPop(key []byte, values [][]byte) (bool, error) {
	lock(LSortedSet, key)
	defer unlock(LSortedSet, key)

	batch := store.NewBatch()
	defer batch.Close()

	for _, value := range values {
		exist, err := store.Get(generateSortedSetScoreKey(key, value))
		if err != nil {
			return false, err
		}
		if exist != nil && len(exist) > 0 {
			existTuple := pb.Tuple{}
			if err := proto.Unmarshal(exist, &existTuple); err != nil {
				return false, err
			}
			batch.Del(generateSortedSetScoreKey(key, value))
			batch.Del(generateSortedSetTupleKey(key, existTuple.Score, value))
		}
	}

	return batch.Commit()
}

func ZRange(key []byte, offset int32, limit uint32) ([]*pb.Tuple, error) {
	index := int32(0)
	res := make([]*pb.Tuple, limit)
	store.Iterate(&engine.PrefixIteratorOption{Prefix: generateSortedSetTupleKeyPrefix(key),
		Offset: offset, Limit: limit},
		func(key []byte, value []byte) {
			tuple := pb.Tuple{}
			_ = proto.Unmarshal(value, &tuple)
			res[index] = &tuple
			index++
		})
	return res[0:index], nil
}

func ZExist(key []byte, values [][]byte) ([]bool, error) {
	res := make([]bool, len(values))
	for i, value := range values {
		exist, err := store.Get(generateSortedSetScoreKey(key, value))
		if err != nil {
			return nil, err
		}
		res[i] = len(exist) > 0
	}
	return res, nil
}

func ZDel(key []byte) (bool, error) {
	lock(LSortedSet, key)
	defer unlock(LSortedSet, key)

	batch := store.NewBatch()
	defer batch.Close()

	store.Iterate(&engine.PrefixIteratorOption{Prefix: generateSortedSetScoreKeyPrefix(key)},
		func(key []byte, value []byte) {
			batch.Del(key)
		})

	store.Iterate(&engine.PrefixIteratorOption{Prefix: generateSortedSetTupleKeyPrefix(key)},
		func(key []byte, value []byte) {
			batch.Del(key)
		})

	return batch.Commit()
}

func ZCount(key []byte) (uint32, error) {
	count := uint32(0)
	store.Iterate(&engine.PrefixIteratorOption{Prefix: generateSortedSetTupleKeyPrefix(key)},
		func(_ []byte, _ []byte) {
			count++
		})
	return count, nil
}

func ZMembers(key []byte) ([]*pb.Tuple, error) {
	index := int32(0)
	res := make([]*pb.Tuple, 0)
	store.Iterate(&engine.PrefixIteratorOption{Prefix: generateSortedSetTupleKeyPrefix(key),
		Offset: 0, Limit: math.MaxInt32},
		func(key []byte, value []byte) {
			// zs/{key}/{score}/{value} -> {tuple}
			tuple := pb.Tuple{}
			_ = proto.Unmarshal(value, &tuple)
			res = append(res, &tuple)
			index++
		})
	return res[0:index], nil
}

func generateSortedSetScoreKey(key []byte, value []byte) []byte {
	return []byte(fmt.Sprintf(sortedSetScoreKeyTemplate, key, value))
}

func generateSortedSetScoreKeyValue(score float64) []byte {
	return []byte(fmt.Sprintf("%32.32f", score))
}

func generateSortedSetScoreKeyPrefix(key []byte) []byte {
	return []byte(fmt.Sprintf(sortedSetScoreKeyPrefixTemplate, key))
}

func generateSortedSetTupleKeyPrefix(key []byte) []byte {
	return []byte(fmt.Sprintf(sortedSetTupleKeyPrefixTemplate, key))
}

func generateSortedSetTupleKey(key []byte, score float64, value []byte) []byte {
	return []byte(fmt.Sprintf(sortedSetTupleKeyTemplate, key, generateSortedSetScoreKeyValue(score), value))
}
