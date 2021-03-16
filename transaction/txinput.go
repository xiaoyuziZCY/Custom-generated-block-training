package transaction

//定义交易输入结构体
type TxInput struct {
	Txid [32]byte
	Vout int//引用自哪个交易输出
	ScriptSig []byte
}
