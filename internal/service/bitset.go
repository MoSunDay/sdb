package service

import (
	"errors"
	"fmt"
	"github.com/tmthrgd/go-bitset"
	"github.com/yemingfeng/sdb/internal/store"
)

var NotFoundBitsetError = errors.New("not found bitset, please create it")
var BitsetExistError = errors.New("bitset exist, please delete it or change other")
var BitsetRangeError = errors.New("bitset out of range, please check it")

const bitsetKeyTemplate = "bs/%s"

func BSCreate(key []byte, size uint32) (bool, error) {
	lock(LBitset, key)
	defer unlock(LBitset, key)

	bitsetKey := generateBitsetKey(key)
	bitsetValue, err := store.Get(bitsetKey)
	if err != nil {
		return false, err
	}
	if bitsetValue != nil && len(bitsetValue) > 0 {
		return false, BitsetExistError
	}

	return store.Set(bitsetKey, bitset.New(uint(size)))
}

func BSDel(key []byte) (bool, error) {
	return store.Del(generateBitsetKey(key))
}

func BSSetRange(key []byte, start uint32, end uint32, value bool) (bool, error) {
	lock(LBitset, key)
	defer unlock(LBitset, key)

	bitsetKey := generateBitsetKey(key)
	bitsetValue, err := store.Get(bitsetKey)
	if err != nil {
		return false, err
	}
	if bitsetValue == nil {
		return false, NotFoundBitsetError
	}
	b := bitset.Bitset(bitsetValue)
	if start > end {
		return false, BitsetRangeError
	}
	if end > uint32(b.Len()) {
		return false, BitsetRangeError
	}
	b.SetRangeTo(uint(start), uint(end), value)
	return store.Set(bitsetKey, b)
}

func BSMSet(key []byte, bits []uint32, value bool) (bool, error) {
	lock(LBitset, key)
	defer unlock(LBitset, key)

	bitsetKey := generateBitsetKey(key)
	bitsetValue, err := store.Get(bitsetKey)
	if err != nil {
		return false, err
	}
	if bitsetValue == nil {
		return false, NotFoundBitsetError
	}
	b := bitset.Bitset(bitsetValue)
	for i := range bits {
		bit := uint(bits[i])
		if bit > b.Len() {
			return false, BitsetRangeError
		}
		b.SetTo(bit, value)
	}
	return store.Set(bitsetKey, b)
}

func BSGetRange(key []byte, start uint32, end uint32) ([]bool, error) {
	bitsetKey := generateBitsetKey(key)
	bitsetValue, err := store.Get(bitsetKey)
	if err != nil {
		return nil, err
	}
	if bitsetValue == nil {
		return nil, NotFoundBitsetError
	}
	b := bitset.Bitset(bitsetValue)
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
	bitsetKey := generateBitsetKey(key)
	bitsetValue, err := store.Get(bitsetKey)
	if err != nil {
		return nil, err
	}
	if bitsetValue == nil {
		return nil, NotFoundBitsetError
	}
	b := bitset.Bitset(bitsetValue)
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
	bitsetKey := generateBitsetKey(key)
	bitsetValue, err := store.Get(bitsetKey)
	if err != nil {
		return 0, err
	}
	if bitsetValue == nil {
		return 0, NotFoundBitsetError
	}
	b := bitset.Bitset(bitsetValue)
	return uint32(b.Count()), nil
}

func BSCountRange(key []byte, start uint32, end uint32) (uint32, error) {
	bitsetKey := generateBitsetKey(key)
	bitsetValue, err := store.Get(bitsetKey)
	if err != nil {
		return 0, err
	}
	if bitsetValue == nil {
		return 0, NotFoundBitsetError
	}
	b := bitset.Bitset(bitsetValue)
	if start > end {
		return 0, BitsetRangeError
	}
	if end > uint32(b.Len()) {
		return 0, BitsetRangeError
	}
	return uint32(b.CountRange(uint(start), uint(end))), nil
}

func generateBitsetKey(key []byte) []byte {
	return []byte(fmt.Sprintf(bitsetKeyTemplate, key))
}
