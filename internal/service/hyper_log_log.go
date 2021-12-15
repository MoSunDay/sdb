package service

import (
	"errors"
	"fmt"
	"github.com/axiomhq/hyperloglog"
	"github.com/yemingfeng/sdb/internal/store"
)

var NotFoundHyperLogLogError = errors.New("not found hyper log log, please create it")
var HyperLogLogExistError = errors.New("hyper log log exist, please delete it or change other")

const hyperLogLogKeyTemplate = "hl/%s"

func HLLCreate(key []byte) (bool, error) {
	lock(LHyperLogLog, key)
	defer unlock(LHyperLogLog, key)

	hyperLogLogKey := generateHyperLogLogKey(key)
	value, e := store.Get(hyperLogLogKey)
	if e != nil {
		return false, e
	}
	if value != nil && len(value) > 0 {
		return false, HyperLogLogExistError
	}

	h := hyperloglog.New16()
	value, e = h.MarshalBinary()
	if e != nil {
		return false, e
	}

	return store.Set(hyperLogLogKey, value)
}

func HLLDel(key []byte) (bool, error) {
	return store.Del(generateHyperLogLogKey(key))
}

func HLLAdd(key []byte, values [][]byte) (bool, error) {
	lock(LHyperLogLog, key)
	defer unlock(LHyperLogLog, key)

	value, err := store.Get(generateHyperLogLogKey(key))
	if err != nil {
		return false, err
	}
	if len(value) == 0 {
		return false, NotFoundHyperLogLogError
	}
	var hll hyperloglog.Sketch
	err = hll.UnmarshalBinary(value)
	if err != nil {
		return false, err
	}

	for _, value := range values {
		hll.Insert(value)
	}

	value, e := hll.MarshalBinary()
	if e != nil {
		return false, e
	}
	return store.Set(generateHyperLogLogKey(key), value)
}

func HLLCount(key []byte) (uint32, error) {
	value, err := store.Get(generateHyperLogLogKey(key))
	if err != nil {
		return 0, err
	}
	if len(value) == 0 {
		return 0, NotFoundHyperLogLogError
	}
	var hll hyperloglog.Sketch
	err = hll.UnmarshalBinary(value)
	if err != nil {
		return 0, err
	}
	return uint32(hll.Estimate()), nil
}

func generateHyperLogLogKey(key []byte) []byte {
	return []byte(fmt.Sprintf(hyperLogLogKeyTemplate, key))
}
