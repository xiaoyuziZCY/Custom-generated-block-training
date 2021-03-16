package consensus

import (
	"Xianfeng/transaction"
	"math/big"
)

type Consensus interface {
	SearchNonce() ([32]byte,int64)
}
//区块的数据标准
type Blockinterface interface {
	Getheight()int64
	Getversion()int64
	Gettimestamp()int64
	Getprehash()[32]byte
	Gettxs() []transaction.Transaction
}
func NewPow(block Blockinterface) Consensus{
	init := big.NewInt(1)
	init.Lsh(init,255-DIFFICULTY)
	return POW{block,init}
}
func NewPos() Consensus{
	proof:= POS{}
	return proof
}