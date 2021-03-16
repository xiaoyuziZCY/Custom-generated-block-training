package transaction

import (
	"Xianfeng/utils"
	"crypto/sha256"
)

const REWARD  =50
//定义交易信息的结构体
type Transaction struct {
	TxHash [32]byte
	Inputs []TxInput
	Outputs []TxOutput
}
//该函数用于构建一个coinbase交易，返回交易结构体实例
func NewCoinbaseTx(address string)(*Transaction,error){
	//txInput:=TxInput{}
	//构造
	txOutput:=TxOutput{
		Value:     REWARD,
		ScriptPub: []byte(address),
	}
	tx:=Transaction{
		Inputs:[]TxInput{},
		Outputs:[]TxOutput{txOutput},
	}
	//序列化
	txBytes,err:=utils.GobEncode(tx)
	if err != nil {
		return nil,err
	}
	//交易哈希
	tx.TxHash=sha256.Sum256(txBytes)
	return &tx,nil
}
//构建一个新的交易
func NewTransaction(from string,to string,value float64)(*Transaction,error){
	//便利区块并找到所有交易
	inputs:=make([]TxInput,0)
	var inputAmout float64
	outputs:=make([]TxOutput,0)
	TxOutput0:=TxOutput{
		Value:value,
		ScriptPub:[]byte(to),
	}
	outputs=append(outputs,TxOutput0)
	if inputAmout-value>0 {
		TxOutput1:=TxOutput{
			Value:     inputAmout-value,
			ScriptPub: []byte(from),
		}
		outputs=append(outputs,TxOutput1)
	}
	tx:=Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}
	txBytes,err:=utils.GobEncode(tx)
	if err != nil {
		return nil,err
	}
	tx.TxHash=sha256.Sum256(txBytes)
	return &tx,nil
}