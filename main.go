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
	fmt.Println("233")
	blockChain,err:=chain.NewBlockChain(engine)
	if err != nil {
		panic(err.Error())
		return
	}
	//fmt.Println(engine)
	cli:=client.Client{blockChain}
	cli.Run()

}