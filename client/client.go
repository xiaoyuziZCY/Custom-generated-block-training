package client

import (
	"Xianfeng/chain"
	"flag"
	"fmt"
	"math/big"
	"os"
)

//客户端（命令行窗口工具），实现与用户进行动态交互
type Client struct {
Chain chain.Blockchain
}
//client的run方法，主要处理逻辑入口
func (client *Client)Run(){
	if len(os.Args)==1 {//用户没输入任何指令
		client.Help()
		return
	}
	command:=os.Args[1]
	switch command {
	case CREATCHAIN:
		fmt.Println("调用创建区块链功能")
		flag.NewFlagSet("createchain",flag.ExitOnError)

	case GENERATEGENESIS:
		fmt.Println("调用生成创世区块的功能")
		generateGensis:=flag.NewFlagSet("generategensis",flag.ExitOnError)
		gensis:=generateGensis.String("gensis","","")
		generateGensis.Parse(os.Args[2:])
		//判断是否已经存在创世区块
		hashBig:=new(big.Int)
		hashBig.SetBytes(client.Chain.LastBlock.Hash[:])
		if hashBig.Cmp(big.NewInt(0))==1 {
			fmt.Println("抱歉，创世区块已经存在，无法写入")
			return
		}
		client.Chain.Creatgenesis([]byte(*gensis))
	case ADDNEWBLOCK:
		fmt.Println("调用生成新区快的功能")
	case GETLASTBLOCK:
		fmt.Println("获取最新区块的功能")
	case GETALLBLOCK:
		fmt.Println("获取所有区块的功能")
	case GETBLOCKCOUNT:
		fmt.Println("获取区块总数")
	case HELP:
		fmt.Println("获取使用说明")
	default:
		fmt.Println("抱歉，目前不支持该命令功能")
		fmt.Println("请你使用help功能参考说明")

	}
}
//该方法向控制太输出项目使用说明
func (client *Client)Help(){


	fmt.Println("----------欢迎进入XIANFENGChain项目-------")
	fmt.Println()
	fmt.Println("以下是系统使用说明")
	fmt.Println("\tgo run main.go command [arguments]")
	fmt.Println()
	fmt.Println("当前支持的功能:")
	fmt.Println()
	fmt.Println("\t"+CREATCHAIN+"    创建一条区块链")
	fmt.Println("\t"+GENERATEGENESIS+" 生成创世区块 可接收一个参数-gensis表示创世区块数据")
	fmt.Println("\t"+ADDNEWBLOCK+"    创建一个新区块 可接收一个参数-data表示区块的数据")
	fmt.Println("\t"+GETALLBLOCK+"      ")
	fmt.Println("\t"+GETLASTBLOCK+"   ")
	fmt.Println("\t"+GETBLOCKCOUNT+"  ")

	fmt.Println()
	fmt.Println("其他命令用法请使用help功能")
	fmt.Println()
}
