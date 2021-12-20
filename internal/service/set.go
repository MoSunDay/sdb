package service

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/store"
)

func generateSetPrefix(key []byte) []byte {
	return []byte(fmt.Sprintf("sk/%s/", key))
}

func SPush(key []byte, values [][]byte) (bool, error) {
	lock(LSet, key)
	defer unlock(LSet, key)

	batch := store.NewBatch()
	defer batch.Close()

	for _, value := range values {
		if _, err := store.NewRow1(generateSetPrefix(key), value, value).Set(batch); err != nil {
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
		if _, err := store.NewRow0(generateSetPrefix(key), value).Del(batch); err != nil {
			return false, err
		}
	}

	return batch.Commit()
}

func SExist(key []byte, values [][]byte) ([]bool, error) {
	res := make([]bool, len(values))
	for i, value := range values {
		exist, err := store.NewRow0(generateSetPrefix(key), value).Exist()
		if err != nil {
			return nil, err
		}
		res[i] = exist
	}
	return res, nil
}

func SDel(key []byte) (bool, error) {
	lock(LSet, key)
	defer unlock(LSet, key)

	batch := store.NewBatch()
	defer batch.Close()

	if _, err := store.Del0(generateSetPrefix(key), batch); err != nil {
		return false, err
	}

	return batch.Commit()
}

func SCount(key []byte) (uint32, error) {
	return store.Count(generateSetPrefix(key))
}

func SMembers(key []byte) ([][]byte, error) {
	rows, err := store.Page0(generateSetPrefix(key))
	if err != nil {
		return nil, err
	}
	res := make([][]byte, len(rows))
	for i := range rows {
		res[i] = rows[i].Value
	}
	return res, nil
}
