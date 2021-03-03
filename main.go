package main

import (
	"Xianfeng/chain"
	"fmt"
)

func main() {
	fmt.Println("正常")
	blockchain:=chain.Creatchainwithgenesis([]byte("hello"))
	//genesis:= chain.CreatGenesisBlock([]byte(string("小鱼子")))
	blockchain.Addnewblock([]byte("block1"))
	blockchain.Addnewblock([]byte("block2"))
	fmt.Println("当前区块个数",len(blockchain.Blocks))
	fmt.Println(blockchain.Blocks[0])
	fmt.Println(blockchain.Blocks[1])
	fmt.Println(blockchain.Blocks[2])
}
