package consensus

import (
	"Xianfeng/consensus/pos"
	"Xianfeng/consensus/pow"
)

type Consensus interface {
	Run() interface{}
}
func NewPow() Consensus{
	proof:=pow.POW{}
	return proof
}
func NewPos() Consensus{
	proof:=pos.POS{}
	return proof
}