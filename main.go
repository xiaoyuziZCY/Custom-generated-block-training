package main

import (
	"Xianfeng/chain"
	"fmt"
)

func main() {
	fmt.Println("正常")
	genesis:= chain.CreatGenesisBlock([]byte(string("小鱼子")))
	fmt.Println("创世",genesis)
	fmt.Println("区块hash值是",genesis.Hash)
	block1 :=chain.CreateBlock(genesis.Height,genesis.Hash,nil)
	fmt.Println(block1)
	block2 := chain.CreateBlock(block1.Height,block1.Hash,nil)
	fmt.Println(block2)
}
