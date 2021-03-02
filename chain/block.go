package chain

import (
	"Xianfeng/consensus"
	"Xianfeng/utils"
	"bytes"
	"crypto/sha256"
	"time"
)

const VERSION  = 2
type Block struct {
	Height int64
	Version int64
	preHash [32]byte
	Hash [32]byte//区块hash
	//默克尔根
	Timestamp int64
	//Difficulty int64
	Nonce int64
	Data []byte//区块体

}
//只有区块才能调用该方法
func (block *Block)SetHash(){
	heightByte,_ := utils.Int2Byte(block.Height)
	versionByte,_:=utils.Int2Byte(block.Version)
	timestampByte,_ :=utils.Int2Byte(block.Timestamp)
	nonceByte,_:=utils.Int2Byte(block.Nonce)
   bk:= bytes.Join([][]byte{heightByte,versionByte,block.preHash[:],timestampByte,nonceByte,block.Data},[]byte{})
   blockhash:= sha256.Sum256(bk)
   block.Hash =blockhash
}
//新区块函数
func CreateBlock(height int64,preHash [32]byte,data []byte)Block  {
	block :=Block{}
	block.Height = height + 1
	block.preHash = preHash
	block.Version = VERSION
	block.Timestamp = time.Now().Unix()
	block.Data = data
	//共识机制切换
	block.SetHash()
	cons := consensus.NewPow()
	cons.Run()

	return block
}
//创世区块函数
func CreatGenesisBlock(data []byte)Block{
	genesis :=Block{}
	genesis.Height = 0
	genesis.preHash = [32]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	genesis.Version = VERSION
	genesis.Timestamp = time.Now().Unix()
	genesis.Nonce = 0
	genesis.Data = data
	genesis.SetHash()
	return genesis
}