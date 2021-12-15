package service

import (
	"encoding/binary"
	"fmt"
	"github.com/yemingfeng/sdb/internal/store"
	"github.com/yemingfeng/sdb/internal/store/engine"
	"github.com/yemingfeng/sdb/pkg/pb"
	"math"
	"strconv"
	"strings"
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
		// get key z/{key}/{value} 获取 score
		existScoreBytes, err := store.Get(generateSortedSetScoreKey(key, tuple.Value))
		if err != nil {
			return false, err
		}

		// remove key zs/{key}/{score}/{value}
		if existScoreBytes != nil {
			existScore := math.Float64frombits(binary.LittleEndian.Uint64(existScoreBytes))
			batch.Del(generateSortedSetTupleKey(key, existScore, tuple.Value))
		}

		// add key z/{key}/{value} -> score
		batch.Set(generateSortedSetScoreKey(key, tuple.Value),
			generateSortedSetScoreKeyValue(tuple.Score))

		// add key zs/{key}/{score}/{value} -> value
		batch.Set(generateSortedSetTupleKey(key, tuple.Score, tuple.Value), tuple.Value)
	}

	return batch.Commit()
}

func ZPop(key []byte, values [][]byte) (bool, error) {
	lock(LSortedSet, key)
	defer unlock(LSortedSet, key)

	batch := store.NewBatch()
	defer batch.Close()

	for _, value := range values {
		existScoreStr, err := store.Get(generateSortedSetScoreKey(key, value))
		if err != nil {
			return false, err
		}
		if existScoreStr != nil {
			batch.Del(generateSortedSetScoreKey(key, value))
			existScore, err := strconv.ParseFloat(string(existScoreStr), 64)
			if err != nil {
				return false, err
			}
			batch.Del(generateSortedSetTupleKey(key, existScore, value))
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
			// zs/{key}/{score}/{value} -> {value}
			infos := strings.Split(string(key), "/")
			scoreStr := infos[2]
			score, _ := strconv.ParseFloat(scoreStr, 64)
			res[index] = &pb.Tuple{Score: score, Value: value}
			index++
		})
	return res[0:index], nil
}

func ZExist(key []byte, values [][]byte) ([]bool, error) {
	res := make([]bool, len(values))
	for i, value := range values {
		scoreStr, err := store.Get(generateSortedSetScoreKey(key, value))
		if err != nil {
			return nil, err
		}
		res[i] = len(scoreStr) > 0
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
			// zs/{key}/{score}/{value} -> {value}
			infos := strings.Split(string(key), "/")
			scoreStr := infos[2]
			score, _ := strconv.ParseFloat(scoreStr, 64)
			res = append(res, &pb.Tuple{Score: score, Value: value})
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
