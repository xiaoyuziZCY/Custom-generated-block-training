package main

import (
	"Xianfeng/chain"
	"Xianfeng/client"
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
	blockChain:=chain.Newblockchain(engine)
	//fmt.Println(engine)
	cli:=client.Client{blockChain}
	cli.Run()

}