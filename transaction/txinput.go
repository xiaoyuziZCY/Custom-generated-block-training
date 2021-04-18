package transaction

import "bytes"

/**
 * 定义交易输入结构体
 */
type TxInput struct {
	TxId [32]byte //标识引用自哪笔交易
	Vout int      //引用自哪个交易输出
	//ScriptSig []byte   //解锁脚本 ：交易签名 原始公钥
	//ScriptSig = Sig + PubK
	Sig  []byte //签名
	Pubk []byte //原始公钥
}

/**
 * 判断某个地址对应的原始公钥是否与input相匹配
 */
func (input *TxInput) CheckPubKWithInput(pubk []byte) bool {
	return bytes.Compare(input.Pubk, pubk) == 0
}

/**
 * 构建一个新的交易输入，并将结构体实例返回
 */
func NewTxInput(txid [32]byte, vout int, pubk []byte) TxInput {
	input := TxInput{
		TxId: txid,
		Vout: vout,
		Pubk: pubk,
	}
	//缺两个字段：sig和Pubk还未赋值
	return input
}