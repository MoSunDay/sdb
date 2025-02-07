package service

import (
	"errors"
	"github.com/tmthrgd/go-bitset"
	"github.com/yemingfeng/sdb/internal/collection"
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
)

var NotFoundBitsetError = errors.New("not found bitset, please create it")
var BitsetExistError = errors.New("bitset exist, please delete it or change other")
var BitsetRangeError = errors.New("bitset out of range, please check it")

var bitsetCollection = collection.NewCollection(pb.DataType_BITSET)

func BSCreate(key []byte, size uint32, batch engine.Batch) error {
	exist, err := bitsetCollection.ExistRowById(key, key)
	if err != nil {
		return err
	}
	if exist {
		return BitsetExistError
	}

	return bitsetCollection.UpsertRow(&collection.Row{
		Key:   key,
		Id:    key,
		Value: bitset.New(uint(size)),
	}, batch)
}

func BSDel(key []byte, batch engine.Batch) error {
	return bitsetCollection.DelRowById(key, key, batch)
}

func BSSetRange(key []byte, start uint32, end uint32, value bool, batch engine.Batch) error {
	row, err := bitsetCollection.GetRowById(key, key)
	if err != nil {
		return err
	}
	if row == nil {
		return NotFoundBitsetError
	}
	b := bitset.Bitset(row.Value)
	if start > end {
		return BitsetRangeError
	}
	if end > uint32(b.Len()) {
		return BitsetRangeError
	}
	b.SetRangeTo(uint(start), uint(end), value)
	return bitsetCollection.UpsertRow(&collection.Row{
		Key:   key,
		Id:    key,
		Value: b,
	}, batch)
}

func BSMSet(key []byte, bits []uint32, value bool, batch engine.Batch) error {
	row, err := bitsetCollection.GetRowById(key, key)
	if err != nil {
		return err
	}
	if row == nil {
		return NotFoundBitsetError
	}
	b := bitset.Bitset(row.Value)
	for i := range bits {
		bit := uint(bits[i])
		if bit > b.Len() {
			return BitsetRangeError
		}
		b.SetTo(bit, value)
	}
	return bitsetCollection.UpsertRow(&collection.Row{
		Key:   key,
		Id:    key,
		Value: b,
	}, batch)
}

func BSGetRange(key []byte, start uint32, end uint32) ([]bool, error) {
	row, err := bitsetCollection.GetRowById(key, key)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, NotFoundBitsetError
	}
	b := bitset.Bitset(row.Value)
	if start > end {
		return nil, BitsetRangeError
	}
	if end > uint32(b.Len()) {
		return nil, BitsetRangeError
	}
	res := make([]bool, end-start)
	for i := start; i < end; i++ {
		res[i-start] = b.IsSet(uint(i))
	}
	return res, nil
}

func BSMGet(key []byte, bits []uint32) ([]bool, error) {
	row, err := bitsetCollection.GetRowById(key, key)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, NotFoundBitsetError
	}
	b := bitset.Bitset(row.Value)
	res := make([]bool, len(bits))
	for i := range bits {
		bit := uint(bits[i])
		if bit > b.Len() {
			return nil, BitsetRangeError
		}
		res[i] = b.IsSet(bit)
	}
	return res, nil
}

func BSCount(key []byte) (uint32, error) {
	row, err := bitsetCollection.GetRowById(key, key)
	if err != nil {
		return 0, err
	}
	if row == nil {
		return 0, NotFoundBitsetError
	}
	b := bitset.Bitset(row.Value)
	return uint32(b.Count()), nil
}

func BSCountRange(key []byte, start uint32, end uint32) (uint32, error) {
	row, err := bitsetCollection.GetRowById(key, key)
	if err != nil {
		return 0, err
	}
	if row == nil {
		return 0, NotFoundBitsetError
	}
	b := bitset.Bitset(row.Value)
	if start > end {
		return 0, BitsetRangeError
	}
	if end > uint32(b.Len()) {
		return 0, BitsetRangeError
	}
	return uint32(b.CountRange(uint(start), uint(end))), nil
}
