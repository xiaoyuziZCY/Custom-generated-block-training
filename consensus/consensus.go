package consensus

import (
	"math/big"
	"Xianfeng/transaction"
)

/**
 * 共识机制的接口标准,用于定义共识方案的接口
 */
type Consensus interface {
	SearchNonce() ([32]byte, int64)
}

/**
 * 区块的数据接口标准
 */
type BlockInterface interface {
	GetHeight() int64
	GetVersion() int64
	GetTimeStamp() int64
	GetPreHash() [32]byte
	GetTxs() []transaction.Transaction
}

func NewProofWork(block BlockInterface) Consensus {
	init := big.NewInt(1)
	init.Lsh(init, 255-DIFFICULTY)
	return ProofWork{block,init}
}
func NewPos() Consensus{
	proof:= POS{}
	return proof
}