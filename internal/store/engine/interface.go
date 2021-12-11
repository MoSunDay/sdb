package engine

type Store interface {
	Set(key []byte, value []byte, sync bool) (bool, error)
	Get(key []byte) ([]byte, error)
	Del(key []byte, sync bool) (bool, error)
	NewBatch() Batch
	Iterate(opt *PrefixIteratorOption, handle func([]byte, []byte))
	Close() error
}

type Batch interface {
	Set(key []byte, value []byte)
	Del(key []byte)
	Commit(sync bool) (bool, error)
	Close()
}

type PrefixIteratorOption struct {
	Prefix []byte

	Offset int
	Limit  int
}
