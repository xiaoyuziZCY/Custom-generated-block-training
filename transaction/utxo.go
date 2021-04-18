package transaction

import (
	"Xianfeng/utils"
	"Xianfeng/wallet"
	"bytes"
)

//未花费的交易输出
type UTXO struct {
	Txid     [32]byte //表明该可花的钱在哪笔交易上
	Vout     int      //表明该可花的钱在该交易的哪个交易输出上
	TxOutput         //用集成的方式表示utxo上有多少可用的钱和该笔钱属于谁
}
func NewUTXO(txid [32]byte,vout int,output TxOutput)UTXO{
	return UTXO{
		Txid:     txid,
		Vout:     vout,
		TxOutput: output,
	}
}
func (utxo *UTXO)IsSpent(spend TxInput)bool{
	equalTxId:=bytes.Compare(utxo.Txid[:],spend.TxId[:])==0
	equalVout:=utxo.Vout==spend.Vout
	pubk:=spend.Pubk
	pubk256:=utils.Hash256(pubk)
	ripemd160:=utils.RipEMd160(pubk256)
	versionPubkHash:=append([]byte(wallet.VERSION),ripemd160...)
	equalLock:=bytes.Compare(utxo.PubKHash,versionPubkHash)==0
	return equalTxId&&equalVout&&equalLock
}
func (utxo *UTXO) EqualUTXO(specified UTXO) bool {
	equalTxId := bytes.Compare(utxo.Txid[:], specified.Txid[:]) == 0
	equalVout := utxo.Vout == specified.Vout
	equalValue := utxo.Value == specified.Value
	equalPubKHash := bytes.Compare(utxo.PubKHash, specified.PubKHash) == 0
	return equalTxId && equalVout && equalValue && equalPubKHash
}