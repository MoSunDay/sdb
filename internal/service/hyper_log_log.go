package service

import (
	"errors"
	"github.com/axiomhq/hyperloglog"
	"github.com/yemingfeng/sdb/internal/collection"
	"github.com/yemingfeng/sdb/internal/pb"
)

var NotFoundHyperLogLogError = errors.New("not found hyper log log, please create it")
var HyperLogLogExistError = errors.New("hyper log log exist, please delete it or change other")

var hyperLogLogCollection = collection.NewCollection(pb.DataType_HYPER_LOG_LOG)

func HLLCreate(key []byte) (bool, error) {
	lock(LHyperLogLog, key)
	defer unlock(LHyperLogLog, key)

	exist, err := hyperLogLogCollection.ExistRowById(key, key)
	if err != nil {
		return false, err
	}
	if exist {
		return false, HyperLogLogExistError
	}

	h := hyperloglog.New16()
	value, err := h.MarshalBinary()
	if err != nil {
		return false, err
	}

	return hyperLogLogCollection.UpsertRowAutoCommit(&collection.Row{
		Key:   key,
		Id:    key,
		Value: value,
	})
}

func HLLDel(key []byte) (bool, error) {
	return hyperLogLogCollection.DelRowByIdAutoCommit(key, key)
}

func HLLAdd(key []byte, values [][]byte) (bool, error) {
	lock(LHyperLogLog, key)
	defer unlock(LHyperLogLog, key)

	row, err := hyperLogLogCollection.GetRowById(key, key)
	if err != nil {
		return false, err
	}
	if row == nil {
		return false, NotFoundHyperLogLogError
	}
	var hll hyperloglog.Sketch
	err = hll.UnmarshalBinary(row.Value)
	if err != nil {
		return false, err
	}

	for _, value := range values {
		hll.Insert(value)
	}

	value, err := hll.MarshalBinary()
	if err != nil {
		return false, err
	}
	return hyperLogLogCollection.UpsertRowAutoCommit(&collection.Row{
		Key:   key,
		Id:    key,
		Value: value,
	})
}

func HLLCount(key []byte) (uint32, error) {
	row, err := hyperLogLogCollection.GetRowById(key, key)
	if err != nil {
		return 0, err
	}
	if row == nil {
		return 0, NotFoundHyperLogLogError
	}
	var hll hyperloglog.Sketch
	err = hll.UnmarshalBinary(row.Value)
	if err != nil {
		return 0, err
	}
	return uint32(hll.Estimate()), nil
}
