package store

import (
	"encoding/json"
	"fmt"
	"github.com/yemingfeng/sdb/internal/store/engine"
)

type Index struct {
	Name  []byte
	Value []byte
}

// Row 每行的数据，包含 prefix id indexes
// 以 List 存储为例子，若 prefix: ll/list，有 4 个元素，分别是：
// { {element: aaa}, {score: 1.1} }
// { {element: bbb}, {score: 2.2} }
// { {element: ccc}, {score: 3.3} }
// { {element: aaa}, {score: 4.4} }
// 则有 4 个 Row 对象，分别是：
// Row1 { prefix: ll/list, id: aaa:1.1，index: [{ name: element, value: aaa }, {name: score, value: 1.1 }], }
// Row2 { prefix: ll/list, id: bbb:2.2，index: [{ name: element, value: bbb }, {name: score, value: 2.2 }] }
// Row3 { prefix: ll/list, id: ccc:.3.3，index: [{ name: element, value: ccc }, {name: score, value: 3.3 }] }
// Row4 { prefix: ll/list, id: aaa:4.4，index: [{ name: element, value: aaa }, {name: score, value: 4.4 }] }
// 写入到存储引擎结构的 rowKey = {prefix}/id
// Row1 的 rowKey 为 ll/list/aaa/1.1 会写入三条数据到存储引擎：
//    ll/list/aaa/1.1 -> { prefix: ll/list, id: aaa/1.1, indexes: [...], value: [....] }
//    {ll/list}/element}/{aaa}/{aaa/1.1} -> ll/list/aaa/1.1
//    {ll/list}/{score}/{1.1}/{aaa/1.1} -> ll/list/aaa/1.1
// 这样就可以同时对 List 中的 element 和 score 检索了。实现了多级索引
type Row struct {
	Prefix  []byte // 前缀标识
	Id      []byte
	Indexes []Index
	Value   []byte
}

func NewRow0(prefix []byte, id []byte) *Row {
	return &Row{Prefix: prefix, Id: id}
}

func NewRow1(prefix []byte, id []byte, value []byte) *Row {
	return &Row{Prefix: prefix, Id: id, Value: value}
}

func NewRow2(prefix []byte, id []byte, indexes []Index, value []byte) *Row {
	return &Row{Prefix: prefix, Id: id, Indexes: indexes, Value: value}
}

func rowKey(prefix []byte, id []byte) []byte {
	return []byte(fmt.Sprintf("%s/%s", prefix, id))
}

func rowPrefixKey(prefix []byte) []byte {
	return []byte(fmt.Sprintf("%s", prefix))
}

func indexKey(prefix []byte, indexName []byte, indexValue []byte, id []byte) []byte {
	return []byte(fmt.Sprintf("%s/%s/%s/%s", prefix, indexName, indexValue, id))
}

func indexPrefixKey(prefix []byte, indexKey []byte) []byte {
	return []byte(fmt.Sprintf("%s/%s/", prefix, indexKey))
}

func (row *Row) Get() (bool, error) {
	exist, err := Get(rowKey(row.Prefix, row.Id))
	if err != nil {
		return false, err
	}
	if (len(exist)) == 0 {
		return true, nil
	}

	existRow := Row{}
	if err := json.Unmarshal(exist, &existRow); err != nil {
		return false, err
	}
	row.Indexes = existRow.Indexes
	row.Value = existRow.Value

	return true, nil
}

func (row *Row) Exist() (bool, error) {
	if _, err := row.Get(); err != nil {
		return false, err
	}
	return len(row.Value) > 0, nil
}

func (row *Row) Set(batch engine.Batch) (bool, error) {
	rowKey := rowKey(row.Prefix, row.Id)

	exist, err := Get(rowKey)
	if err != nil {
		return false, err
	}
	// 得先删除之前的数据
	if exist != nil && len(exist) > 0 {
		if _, err := batch.Del(rowKey); err != nil {
			return false, err
		}
		existRow := Row{}
		if err := json.Unmarshal(exist, &existRow); err != nil {
			return false, err
		}
		if existRow.Indexes != nil {
			for i := range existRow.Indexes {
				existIndex := existRow.Indexes[i]
				if _, err := batch.Del(indexKey(row.Prefix, existIndex.Name, existIndex.Value, row.Id)); err != nil {
					return false, err
				}
			}
		}
	}

	rawValue, err := json.Marshal(row)
	if err != nil {
		return false, err
	}

	if _, err := batch.Set(rowKey, rawValue); err != nil {
		return false, err
	}
	if row.Indexes != nil {
		for i := range row.Indexes {
			index := row.Indexes[i]
			if _, err := batch.Set(indexKey(row.Prefix, index.Name, index.Value, row.Id), rowKey); err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func (row *Row) Del(batch engine.Batch) (bool, error) {
	if _, err := batch.Del(rowKey(row.Prefix, row.Id)); err != nil {
		return false, err
	}
	for i := range row.Indexes {
		index := row.Indexes[i]
		if _, err := batch.Del(indexKey(row.Prefix, index.Name, index.Value, row.Id)); err != nil {
			return false, err
		}
	}
	return true, nil
}
