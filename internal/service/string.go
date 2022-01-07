package service

import (
	"github.com/yemingfeng/sdb/internal/collection"
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
	"strconv"
)

var stringCollection = collection.NewCollection(pb.DataType_STRING)

func Set(key []byte, value []byte, batch engine.Batch) error {
	return stringCollection.UpsertRow(&collection.Row{
		Key:   key,
		Id:    key,
		Value: value}, batch)
}

func MSet(keys [][]byte, values [][]byte, batch engine.Batch) error {
	for i := range keys {
		if err := stringCollection.UpsertRow(&collection.Row{
			Key:   keys[i],
			Id:    keys[i],
			Value: values[i],
		}, batch); err != nil {
			return err
		}
	}
	return nil
}

func SetNX(key []byte, value []byte, batch engine.Batch) error {
	exist, err := stringCollection.ExistRowById(key, key)
	if err != nil {
		return err
	}
	if exist {
		return err
	}
	return stringCollection.UpsertRow(&collection.Row{
		Key:   key,
		Id:    key,
		Value: value}, batch)
}

func Get(key []byte) ([]byte, error) {
	row, err := stringCollection.GetRowById(key, key)
	if err != nil || row == nil {
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

func Del(key []byte, batch engine.Batch) error {
	return stringCollection.Del(key, batch)
}

func Incr(key []byte, delta int32, batch engine.Batch) error {
	row, err := stringCollection.GetRowById(key, key)
	if err != nil {
		return err
	}
	var valueInt = 0
	if row != nil {
		valueInt, err = strconv.Atoi(string(row.Value))
		if err != nil {
			return err
		}
	}
	valueInt = valueInt + int(delta)

	return stringCollection.UpsertRow(&collection.Row{
		Key:   key,
		Id:    key,
		Value: []byte(strconv.Itoa(valueInt))}, batch)
}
