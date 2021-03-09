package chain

type ChainIterator interface {
	HasNext() bool
	Next() Block
}
