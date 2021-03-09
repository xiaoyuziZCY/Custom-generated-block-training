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
	err = blockChain.Addnewblock([]byte("先锋小镇"))
	if err != nil {
		fmt.Println("错误信息1：", err.Error())
		return
	}
	for blockChain.HasNext() {
		block:=blockChain.Next()
		fmt.Printf("区块%d",block.Height)
		fmt.Printf("区块hash：%v",block.Hash)
		fmt.Printf("前区块hash：%v",block.PreHash)
		fmt.Printf("区块数据%s",block.Data)
		fmt.Println()
	}
}