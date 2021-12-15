package service

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/store"
	"strconv"
)

const stringKeyPrefixTemplate = "s/%s"

func Set(key []byte, value []byte) (bool, error) {
	lock(LString, key)
	defer unlock(LString, key)

	return store.Set(generateStringKey(key), value)
}

func MSet(keys [][]byte, values [][]byte) (bool, error) {
	batch := store.NewBatch()
	defer batch.Close()

	for i := range keys {
		batch.Set(generateStringKey(keys[i]), values[i])
	}
	return batch.Commit()
}

func SetNX(key []byte, value []byte) (bool, error) {
	lock(LString, key)
	defer unlock(LString, key)

	storeKey := generateStringKey(key)
	oldValue, err := store.Get(storeKey)
	if err != nil {
		return false, err
	}
	if len(oldValue) > 0 {
		return false, err
	}
	return store.Set(storeKey, value)
}

func SetGet(key []byte, value []byte) (bool, []byte, error) {
	lock(LString, key)
	defer unlock(LString, key)

	storeKey := generateStringKey(key)
	oldValue, err := store.Get(storeKey)
	if err != nil {
		return false, nil, err
	}
	res, err := store.Set(storeKey, value)
	return res, oldValue, err
}

func Get(key []byte) ([]byte, error) {
	return store.Get(generateStringKey(key))
}

func MGet(keys [][]byte) ([][]byte, error) {
	values := make([][]byte, len(keys))
	for i := range keys {
		value, err := store.Get(generateStringKey(keys[i]))
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	return values, nil
}

func Del(key []byte) (bool, error) {
	return store.Del(generateStringKey(key))
}

func Incr(key []byte, delta int32) (bool, error) {
	lock(LString, key)
	defer unlock(LString, key)

	value, err := store.Get(key)
	if err != nil {
		return false, err
	}
	var valueInt = 0
	if len(value) > 0 {
		valueInt, err = strconv.Atoi(string(value))
		if err != nil {
			return false, err
		}
	}
	valueInt = valueInt + int(delta)

	return store.Set(key, []byte(strconv.Itoa(valueInt)))
}

func generateStringKey(key []byte) []byte {
	return []byte(fmt.Sprintf(stringKeyPrefixTemplate, key))
}
