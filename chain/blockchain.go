package chain


//用于存储在内存里产生的区块
type Blockchain struct {
	Blocks []Block
}
//func Newblockchain()Blockchain{
//	return Blockchain{}
//}
//创建一个区块链，携带创世区块
func Creatchainwithgenesis(genesisdata []byte)Blockchain{
	geneis:=CreatGenesisBlock(genesisdata)
	blocks:=make([]Block,0)
	blocks=append(blocks,geneis)
	return Blockchain{blocks}
}
func (chain *Blockchain) Addnewblock(data []byte){
	//找到切片最后一个数值，为最新区块
	lastblock:=chain.Blocks[len(chain.Blocks)-1]
	//根据最后一个区块产生新区快
	newblock:=CreateBlock(lastblock.Height,lastblock.Hash,data)
	//把新区块放入到切片里
	chain.Blocks=append(chain.Blocks,newblock)

}
