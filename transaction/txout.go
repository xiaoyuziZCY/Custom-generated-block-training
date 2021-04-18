package transaction

import (
	"Xianfeng/utils"
	"bytes"
	"fmt"
)

/**
 * 定义交易输出结构体
 */
type TxOutput struct {
	Value float64 //锁定的币的数量
	//ScriptPub []byte  //锁定的脚本: 指令操作符 + PubHash + 指令操作符
	PubKHash []byte //公钥hash
}

//地址add1
//私钥 -> 公钥 -> hash 160 +version 2次hash check  pinjie base58 -> add1
/**
 * 构建一个新的交易输出，锁一定数量的钱到某个交易输出上，并将该交易输出返回
 */
func Lock2Address(value float64, add string) TxOutput {
	fmt.Println(add)
	reAdd := utils.Decode(add)
	pubHash := reAdd[:len(reAdd)-4]
	output := TxOutput{
		Value:   value,
		PubKHash: pubHash,
	}
	return output
}

/**
 *  该方法用于验证某个交易输出是否是属于某个地址的收入
 */
func (output *TxOutput) VertifyOutputWithAddress(addr string) bool {
	reAdd := utils.Decode(addr)
	pubHash := reAdd[:len(reAdd)-4]
	return bytes.Compare(output.PubKHash, pubHash) == 0
}
