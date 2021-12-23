package service

import (
	"github.com/yemingfeng/sdb/internal/collection"
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
	"math"
)

var mapCollection = collection.NewCollection(pb.DataType_MAP)

func MPush(key []byte, pairs []*pb.Pair) (bool, error) {
	lock(LMap, key)
	defer unlock(LMap, key)

	return mapCollection.Batch(func(batch engine.Batch) error {
		for i := range pairs {
			if _, err := mapCollection.UpsertRow(&collection.Row{
				Key:   key,
				Id:    pairs[i].Key,
				Value: pairs[i].Value,
			}, batch); err != nil {
				return err
			}
		}
		return nil
	})
}

func MPop(key []byte, keys [][]byte) (bool, error) {
	lock(LMap, key)
	defer unlock(LMap, key)

	return mapCollection.Batch(func(batch engine.Batch) error {
		for i := range keys {
			if _, err := mapCollection.DelRowById(key, keys[i], batch); err != nil {
				return err
			}
		}
		return nil
	})
}

func MExist(key []byte, keys [][]byte) ([]bool, error) {
	res := make([]bool, len(keys))
	for i := range keys {
		exist, err := mapCollection.ExistRowById(key, keys[i])
		if err != nil {
			return nil, err
		}
		res[i] = exist
	}
	return res, nil
}

func MDel(key []byte) (bool, error) {
	lock(LMap, key)
	defer unlock(LMap, key)

	return mapCollection.DelAutoCommit(key)
}

func MCount(key []byte) (uint32, error) {
	return mapCollection.Count(key)
}

func MMembers(key []byte) ([]*pb.Pair, error) {
	rows, err := mapCollection.Page(key, 0, math.MaxUint32)
	if err != nil {
		return nil, err
	}

	res := make([]*pb.Pair, len(rows))
	for i := range rows {
		res[i] = &pb.Pair{Key: rows[i].Id, Value: rows[i].Value}
	}

	return res, nil
}
