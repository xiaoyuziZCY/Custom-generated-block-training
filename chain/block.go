package chain

import "time"

const VERSION  = 2
type Block struct {
	Height int64
	Version int64
	preHash [32]byte
	//默克尔根
	Timestamp int64
	//Difficulty int64
	Nonce int64
	Data []byte//区块体

}

func CreateBlock(height int64,preHash [32]byte,Timestamp int64,Nonce int64)Block  {
	block :=Block{}
	block.Height = height + 1
	block.preHash = preHash
	block.Version = VERSION
	block.Timestamp = Timestamp
	block.Nonce = Nonce
	block.Data = nil
	return block
}
func CreatGenesisBlock(data []byte)Block{
	genesis :=Block{}
	genesis.Height = 0
	genesis.preHash = [32]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	genesis.Version = VERSION
	genesis.Timestamp = time.Now().Unix()
	genesis.Nonce = 0
	genesis.Data = data
	return genesis
}