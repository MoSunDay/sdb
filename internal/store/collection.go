package store

import (
	"encoding/json"
	"github.com/yemingfeng/sdb/internal/store/engine"
	"math"
)

func Page0(prefix []byte) ([]*Row, error) {
	return Page1(prefix, 0, math.MaxUint32)
}

func Page1(prefix []byte, offset int32, limit uint32) ([]*Row, error) {
	rows := make([]*Row, 0)
	if err := Iterate1(rowPrefixKey(prefix), offset, limit,
		func(_ []byte, rawRow []byte) error {
			row := Row{}
			if err := json.Unmarshal(rawRow, &row); err != nil {
				return err
			}
			rows = append(rows, &row)
			return nil
		}); err != nil {
		return nil, err
	}
	return rows, nil
}

func IndexPage0(prefix []byte, index []byte) ([]*Row, error) {
	return IndexPage1(prefix, index, 0, math.MaxUint32)
}

func IndexPage1(prefix []byte, index []byte, offset int32, limit uint32) ([]*Row, error) {
	rows := make([]*Row, 0)
	if err := Iterate1(indexPrefixKey(prefix, index), offset, limit,
		func(_ []byte, rowKey []byte) error {
			rawValue, err := Get(rowKey)
			if err != nil {
				return err
			}
			row := Row{}
			if err := json.Unmarshal(rawValue, &row); err != nil {
				return err
			}
			rows = append(rows, &row)
			return nil
		}); err != nil {
		return nil, err
	}
	return rows, nil
}

func Count(prefix []byte) (uint32, error) {
	count := uint32(0)
	if err := Iterate0(rowPrefixKey(prefix), func(rowKey []byte, rawRow []byte) error {
		count += 1
		return nil
	}); err != nil {
		return 0, err
	}
	return count, nil
}

func Del0(prefix []byte, batch engine.Batch) (bool, error) {
	if err := Iterate0(rowPrefixKey(prefix), func(_ []byte, rawRow []byte) error {
		row := Row{}
		if err := json.Unmarshal(rawRow, &row); err != nil {
			return err
		}
		if _, err := row.Del(batch); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return false, err
	}
	return true, nil
}
