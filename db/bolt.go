package db

import (
	"Xianfeng/chain"
	"fmt"
	"github.com/boltdb/bolt"
)

type DBEngine struct {
	DB *bolt.DB
}

//区块存到文件里
func (engine DBEngine)SaveBlock2DB(block chain.Block){
	fmt.Println("将区块存到db中去")
}

//从文件里恢复区块
func (engine DBEngine)GetBlockFromDB(hahs [32]byte)chain.Block{
	fmt.Println("从db中获取特定区块")
}
