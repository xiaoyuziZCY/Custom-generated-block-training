package chain

import (
	"Xianfeng/consensus"
	"bytes"
	"encoding/gob"
	"time"
)

const VERSION  = 2
type Block struct {
	Height  int64
	Version int64
	PreHash [32]byte
	Hash    [32]byte//区块hash
	//默克尔根
	Timestamp int64
	//Difficulty int64
	Nonce int64
	Data []byte//区块体

}
//只有区块才能调用该方法
//func (block *Block)SetHash(){
//	heightByte,_ := utils.Int2Byte(block.Height)
//	versionByte,_:=utils.Int2Byte(block.Version)
//	timestampByte,_ :=utils.Int2Byte(block.Timestamp)
//	nonceByte,_:=utils.Int2Byte(block.Nonce)
//   bk:= bytes.Join([][]byte{heightByte,versionByte,block.preHash[:],timestampByte,nonceByte,block.Data},[]byte{})
//   blockhash:= sha256.Sum256(bk)
//   block.Hash =blockhash
//}
//区块的序列化为[]byte类型
func (block *Block)Serialize()([]byte,error){
	buffer :=new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(&block)
	return buffer.Bytes(),err
}
//区块反序列化，传入[]byte，返回block
func Deserialize(data []byte)(Block,error){
	var block Block
	reader :=new(bytes.Reader)
	reader.Read(data)
	decoder := gob.NewDecoder(reader)
	err :=decoder.Decode(&block)
	return block,err
}
//新区块函数
func CreateBlock(height int64,preHash [32]byte,data []byte)Block  {
	block :=Block{}
	block.Height = height + 1
	block.PreHash = preHash
	block.Version = VERSION
	block.Timestamp = time.Now().Unix()
	block.Data = data
	//共识机制切换
	//block.SetHash()
	cons := consensus.NewPow(block)
	hash,nonce :=cons.Run()
	block.Hash =hash
	block.Nonce = nonce


	return block
}
//创世区块函数
func CreatGenesisBlock(data []byte)Block{
	genesis :=Block{}
	genesis.Height = 0
	genesis.PreHash = [32]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	genesis.Version = VERSION
	genesis.Timestamp = time.Now().Unix()
	genesis.Nonce = 0
	genesis.Data = data
	proof :=consensus.NewPow(genesis)
	hash,nonce :=proof.Run()
	genesis.Hash = hash
	genesis.Nonce=nonce
	return genesis
}
//下面的方法是实现blockinterface的方法主要目的是解决循环导包问题
func(block Block)Getheight()int64{
	return block.Height
}
func(block Block)Getversion()int64{
	return block.Version
}
func(block Block)Gettimestamp()int64{
	return block.Timestamp
}
func(block Block)Getprehash()[32]byte{
	return block.PreHash
}
func(block Block)Getdata()[]byte{
	return block.Data
}