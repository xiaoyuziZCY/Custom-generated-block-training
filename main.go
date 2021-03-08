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
		panic(err.Error())
	}
	blockChain := chain.Newblockchain(engine)
	blockChain.Creatgenesis([]byte("HELLO WORLD"))
	err = blockChain.Addnewblock([]byte("hello"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	allBlock, err := blockChain.GetAllblocks()
	if err !=nil {
		fmt.Println(err.Error())
		return
	}
	for _,block:=range allBlock{
fmt.Println(block)
	}
}