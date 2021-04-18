package transaction

import (
	"Xianfeng/utils"
	"Xianfeng/wallet"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/ecdsa"
	"errors"
	"crypto/rand"
	"fmt"
	"time"
)

const REWARD = 50

/**
 * 定义交易的结构体
 */
type Transaction struct {
	TxHash  [32]byte   //交易的唯一标识
	Inputs  []TxInput  //交易的交易输入
	Outputs []TxOutput //交易的交易输出
	LockedTime int64
}

/**
 * 该函数用于构建一个coinbase交易，返回一个交易的结构体实例
 */
func NewCoinbaseTx(address string) (*Transaction, error) {
	//txInput 为空

	//构造txOutput
	txOutput := Lock2Address(REWARD, address)
	fmt.Println(txOutput)
	tx := Transaction{
		Inputs:  []TxInput{},
		Outputs: []TxOutput{txOutput},
		LockedTime: time.Now().Unix(),
	}

	//序列化
	txBytes, err := utils.GobEncode(tx)
	if err != nil {
		return nil, err
	}
	//交易哈希计算，并赋值给TxHash字段
	tx.TxHash = sha256.Sum256(txBytes)
	return &tx, nil
}

/**
 * 构建一笔新的交易
 */
func NewTransaction(spent []UTXO, from string, pubk []byte, to string, value float64) (*Transaction, error) {
	//交易输入的容器切片
	txInputs := make([]TxInput, 0)
	var inputAmount float64
	for _, utxo := range spent {
		inputAmount += utxo.Value
		input := NewTxInput(utxo.Txid, utxo.Vout, pubk)
		//把构建好的交易输入放入到交易输入容器中
		txInputs = append(txInputs, input)
	}

	//交易输出的容器切片
	//A->B 10 至多会产生两个交易输出
	txOutputs := make([]TxOutput, 0)

	//第一个交易输出：对应转账接收者的输出
	txOutput0 := Lock2Address(value, to)
	txOutputs = append(txOutputs, txOutput0)

	//还有可能产生找零的一个输出：交易发起者给的钱比要转账的钱多
	if inputAmount-value > 0 { //需要找零给转账发起人
		txOutput1 := Lock2Address(inputAmount-value, from)
		txOutputs = append(txOutputs, txOutput1)
	}

	//构建交易
	tx := Transaction{
		Inputs:  txInputs,
		Outputs: txOutputs,
		LockedTime: time.Now().Unix(),
	}
	//序列化
	txBytes, err := utils.GobEncode(tx)
	if err != nil {
		return nil, err
	}
	tx.TxHash = sha256.Sum256(txBytes)

	return &tx, nil
}

/**
 * 使用私钥对某个交易进行交易的签名
 */
func (tx *Transaction) Sign(private *ecdsa.PrivateKey, utxos []UTXO) (error) {

	//交易的交易输入的个数与utxo的个数需要一致
	if tx.IsCoinbaseTransaction() {
		return errors.New("签名遇到错误，请重试")
	}

	if len(tx.Inputs)!=len(utxos) {
		return errors.New("签名遇到错误")
	}
	txCopy:=CopyTx(*tx)

	//tx: 包含多个交易输入TxInput
	for i := 0; i < len(txCopy.Inputs); i++ {
		//input := txCopy.Inputs[i] //当前遍历到的第几个交易输入

		utxo := utxos[i] //当前遍历到的第几个utxo
		//scirptPub := utxo.PubHash //获得当前遍历到的utxo的锁定脚本中的公钥哈希
		txCopy.Inputs[i].Pubk = utxo.PubKHash

		txHash, err := txCopy.CalculateTxHash()
		if err != nil {
			return err
		}

		r, s, err := ecdsa.Sign(rand.Reader, private, txHash)
		if err != nil {
			return err
		}
		txCopy.Inputs[i].Pubk=nil
		sigBytes := append(r.Bytes(), s.Bytes()...)
		tx.Inputs[i].Sig = sigBytes
	}
	return nil
}
func (tx *Transaction)VerifySign(utxos []UTXO)(bool,error){
	if tx.IsCoinbaseTransaction() {
		return true,nil
	}
	if len(tx.Inputs)!=len(utxos) {
		fmt.Print(len(tx.Inputs),len(utxos))
		return false,errors.New("验签错误")
	}
	txCopy:=CopyTx(*tx)
	var result bool
	for index,input :=range txCopy.Inputs {
		pubk:=input.Pubk
		sigBytes:=tx.Inputs[index].Sig
		txCopy.Inputs[index].Sig=nil
		txCopy.Inputs[index].Pubk=utxos[index].PubKHash
		txCopyHash,err:=txCopy.CalculateTxHash()
		if err != nil {
			return false,err
		}
		txCopy.Inputs[index].Pubk=nil
		pub:=wallet.GetPublicKeyWithBytes(elliptic.P256(),pubk)
		r,s:=wallet.RestoreSignature(sigBytes)
		fmt.Printf("签名：%x,%x\n",r.Bytes(),s.Bytes())
		result=ecdsa.Verify(&pub,txCopyHash,r,s)
		if !result {
			return result,errors.New("签名验证失败")
		}
	}
	return result,nil
}

/**
 * 拷贝交易实例
 */
func CopyTx(tx Transaction) Transaction {
	newTx := Transaction{}
	newTx.TxHash=tx.TxHash
	inputs:=make([]TxInput,0)
	for _,input:=range tx.Inputs{
	TxInput:=TxInput{
		TxId: input.TxId,
		Vout: input.Vout,
		Sig:  nil,
		Pubk: input.Pubk,
	}
		inputs=append(inputs,TxInput)
	}
	newTx.Inputs=inputs
	outputs:=make([]TxOutput,0)
	for _,output:=range tx.Outputs{
		TxOutput:=TxOutput{
			Value:    output.Value,
			PubKHash: output.PubKHash,
		}
		outputs=append(outputs,TxOutput)
	}
	newTx.Outputs=outputs
	return newTx
}
func (tx *Transaction)IsCoinbaseTransaction()bool  {
	return len(tx.Inputs)==0 &&len(tx.Outputs)==1
}

/**
 * 计算交易的哈希，并将哈希值返回
 */
func (tx *Transaction) CalculateTxHash() ([]byte, error) {
	txBytes, err := tx.Serialize()
	if err != nil {
		return nil, err
	}
	return utils.Hash256(txBytes), nil
}

/**
 * 交易的序列化
 */
func (tx *Transaction) Serialize() ([]byte, error) {
	return utils.GobEncode(tx)
}