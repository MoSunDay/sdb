package engine

type Store interface {
	Set(key []byte, value []byte) (bool, error)
	Get(key []byte) ([]byte, error)
	Del(key []byte) (bool, error)
	NewBatch() Batch
	Iterate(opt *PrefixIteratorOption, handle func([]byte, []byte))
	Close() error
}

type Batch interface {
	Set(key []byte, value []byte)
	Del(key []byte)
	Commit() (bool, error)
	Close()
}

type PrefixIteratorOption struct {
	Prefix []byte

	Offset int32
	Limit  uint32
}
