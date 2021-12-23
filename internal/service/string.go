package service

import (
	"github.com/yemingfeng/sdb/internal/collection"
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
	"strconv"
)

var stringCollection = collection.NewCollection(pb.DataType_STRING)

func Set(key []byte, value []byte) (bool, error) {
	return stringCollection.UpsertRowAutoCommit(&collection.Row{
		Key:   key,
		Id:    key,
		Value: value})
}

func MSet(keys [][]byte, values [][]byte) (bool, error) {
	return stringCollection.Batch(func(batch engine.Batch) error {
		for i := range keys {
			if _, err := stringCollection.UpsertRow(&collection.Row{
				Key:   keys[i],
				Id:    keys[i],
				Value: values[i],
			}, batch); err != nil {
				return err
			}
		}
		return nil
	})
}

func SetNX(key []byte, value []byte) (bool, error) {
	lock(LString, key)
	defer unlock(LString, key)

	exist, err := stringCollection.ExistRowById(key, key)
	if err != nil {
		return false, err
	}
	if exist {
		return false, err
	}
	return stringCollection.UpsertRowAutoCommit(&collection.Row{
		Key:   key,
		Id:    key,
		Value: value})
}

func SetGet(key []byte, value []byte) (bool, []byte, error) {
	lock(LString, key)
	defer unlock(LString, key)

	oldRow, err := stringCollection.GetRowById(key, key)
	if err != nil {
		return false, nil, err
	}
	res, err := stringCollection.UpsertRowAutoCommit(&collection.Row{
		Key:   key,
		Id:    key,
		Value: value})
	return res, oldRow.Value, err
}

func Get(key []byte) ([]byte, error) {
	row, err := stringCollection.GetRowById(key, key)
	if err != nil {
		return nil, err
	}
	return row.Value, nil
}

func MGet(keys [][]byte) ([][]byte, error) {
	values := make([][]byte, len(keys))
	for i := range keys {
		row, err := stringCollection.GetRowById(keys[i], keys[i])
		if err != nil {
			return nil, err
		}
		if row != nil {
			values[i] = row.Value
		}
	}
	return values, nil
}

func Del(key []byte) (bool, error) {
	return stringCollection.DelAutoCommit(key)
}

func Incr(key []byte, delta int32) (bool, error) {
	lock(LString, key)
	defer unlock(LString, key)

	row, err := stringCollection.GetRowById(key, key)
	if err != nil {
		return false, err
	}
	var valueInt = 0
	if row != nil {
		valueInt, err = strconv.Atoi(string(row.Value))
		if err != nil {
			return false, err
		}
	}
	valueInt = valueInt + int(delta)

	return stringCollection.UpsertRowAutoCommit(&collection.Row{
		Key:   key,
		Id:    key,
		Value: []byte(strconv.Itoa(valueInt))})
}
