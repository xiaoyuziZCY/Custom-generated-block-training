package consensus

import (
	"Xianfeng/utils"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const DIFFICULTY  =16
type POW struct {
		Block  Blockinterface
		Target *big.Int
}

func (pow POW) SearchNonce()([32]byte,int64){
	fmt.Println("已为pow算法机制")
	var nonce int64
	nonce = 0
	hashbig:=new(big.Int)
	for {
		hash := Parepreblock(pow.Block,nonce)
		target :=pow.Target
		hashbig= hashbig.SetBytes(hash[:])
		result :=hashbig.Cmp(target)
		//result := bytes.Compare(hash[:],target.Bytes())
		if result ==-1 {
			return hash,nonce
		}
		nonce++
	}
}
func Parepreblock(block Blockinterface,nonce int64)[32]byte{
	heightByte,_ := utils.Int2Byte(block.Getheight())
	versionByte,_:=utils.Int2Byte(block.Getversion())
	timestampByte,_ :=utils.Int2Byte(block.Gettimestamp())
	nonceByte, _ := utils.Int2Byte(nonce)
	prehash:=block.Getprehash()
	txs := block.Gettxs()

	txsBytes:=make([]byte,0)
		for _,tx:=range txs{
			txBytes,err:=utils.GobEncode(tx)
			if err!=nil {
				break
			}
			//buff.Bytes()
			txsBytes=append(txsBytes,txBytes...)
		}

	bk := bytes.Join([][]byte{
		heightByte,
		versionByte,
		prehash[:],
		timestampByte,
		nonceByte,
		txsBytes}, []byte{})
	return sha256.Sum256(bk)
}