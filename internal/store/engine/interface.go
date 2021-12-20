package engine

type Store interface {
	Get(key []byte) ([]byte, error)
	NewBatch() Batch
	Iterate(opt *PrefixIteratorOption, handle func([]byte, []byte) error) error
	Close() error
}

type Batch interface {
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte) (bool, error)
	Del(key []byte) (bool, error)
	Commit() (bool, error)
	Close()
}

type PrefixIteratorOption struct {
	Prefix []byte

	Offset int32
	Limit  uint32
}
