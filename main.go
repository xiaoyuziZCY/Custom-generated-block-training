package main

import (
	"Xianfeng/chain"
	"fmt"
	"github.com/boltdb/bolt"
)

const DBFILE  ="xianfeng03.db"
func main() {
	fmt.Println("正常")
	engine, err := bolt.Open(DBFILE, 0600, nil)
	if err != nil {
		fmt.Println(0)
		panic(err.Error())
	}
	blockChain := chain.Newblockchain(engine)
	blockChain.Creatgenesis([]byte("HELLO WORLD"))
	lastblock := blockChain.GetLastBlock()
	fmt.Println(lastblock)
	allBlock, err := blockChain.GetAllblocks()
	if err !=nil {
		fmt.Println("pwo...",err.Error())
		return
	}
	fmt.Println(3)
	for _,block:=range allBlock{
fmt.Println(block)
	}
}