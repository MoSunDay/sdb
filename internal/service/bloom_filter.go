package service

import (
	"errors"
	"fmt"
	"github.com/devopsfaith/bloomfilter"
	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
	"github.com/yemingfeng/sdb/internal/store"
)

var NotFoundBloomFilterError = errors.New("not found bloom filter, please create it")
var BloomFilterExistError = errors.New("bloom filter exist, please delete it or change other")

const bloomFilterKeyTemplate = "bf/%s"

func BFCreate(key []byte, n uint32, p float64) (bool, error) {
	lock(LBloomFilter, key)
	defer unlock(LBloomFilter, key)

	bloomFilterKey := generateBloomFilterKey(key)
	value, err := store.Get(bloomFilterKey)
	if err != nil {
		return false, err
	}
	if value != nil && len(value) > 0 {
		return false, BloomFilterExistError
	}
	bloomFilter := baseBloomfilter.New(
		bloomfilter.Config{N: uint(n), P: p, HashName: bloomfilter.HASHER_DEFAULT})

	value, err = bloomFilter.MarshalBinary()
	if err != nil {
		return false, nil
	}

	return store.Set(bloomFilterKey, value)
}

func BFDel(key []byte) (bool, error) {
	return store.Del(generateBloomFilterKey(key))
}

func BFAdd(key []byte, values [][]byte) (bool, error) {
	lock(LBloomFilter, key)
	defer unlock(LBloomFilter, key)

	bloomFilterKey := generateBloomFilterKey(key)
	value, err := store.Get(bloomFilterKey)
	if err != nil {
		return false, err
	}
	if len(value) == 0 {
		return false, NotFoundBloomFilterError
	}

	bloomFilter := &baseBloomfilter.Bloomfilter{}
	if err = bloomFilter.UnmarshalBinary(value); err != nil {
		return false, err
	}

	for _, value := range values {
		bloomFilter.Add(value)
	}

	value, err = bloomFilter.MarshalBinary()
	return store.Set(bloomFilterKey, value)
}

func BFExist(key []byte, values [][]byte) ([]bool, error) {
	bloomFilterKey := generateBloomFilterKey(key)
	value, err := store.Get(bloomFilterKey)
	if err != nil {
		return nil, err
	}
	if len(value) == 0 {
		return nil, NotFoundBloomFilterError
	}
	bloomFilter := &baseBloomfilter.Bloomfilter{}
	err = bloomFilter.UnmarshalBinary(value)
	if err != nil {
		return nil, err
	}

	res := make([]bool, len(values))
	for i, value := range values {
		res[i] = bloomFilter.Check(value)
	}

	return res, nil
}

func generateBloomFilterKey(key []byte) []byte {
	return []byte(fmt.Sprintf(bloomFilterKeyTemplate, key))
}
