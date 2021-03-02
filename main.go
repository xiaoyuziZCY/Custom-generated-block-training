package main

import (
	"Xianfeng/chain"
	"fmt"
	"time"
)

func main() {
	fmt.Println("正常")
	block :=chain.CreateBlock(0,[32]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},time.Now().Unix(),0)

	fmt.Println(block)
	genesis:= chain.CreatGenesisBlock([]byte(string("小鱼子")))
	fmt.Println(genesis)
}
