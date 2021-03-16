package client

import (
	"Xianfeng/chain"
	"Xianfeng/transaction"
	"Xianfeng/utils"
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
		client.GenerateGensis()
	case SENDTRASACtION:
		client.SendTransaction()
	case GETLASTBLOCK:
		client.GetLastBlock()
	case GETALLBLOCKS:
		client.GetAllBlocks()
	case GETBLOCKCOUNT:
		client.GetBlockCount()
	case HELP:
		fmt.Println("获取使用说明")
	default:
		fmt.Println("抱歉，目前不支持该命令功能")
		fmt.Println("请你使用help功能参考说明")

	}
}
func (client *Client)GenerateGensis(){
	fmt.Println("调用生成创世区块的功能")
	generateGensis:=flag.NewFlagSet("generategensis",flag.ExitOnError)
	address:=generateGensis.String("address","","")
	generateGensis.Parse(os.Args[2:])
	//判断是否已经存在创世区块
	hashBig:=new(big.Int)
	hashBig.SetBytes(client.Chain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0))==1 {
		fmt.Println("抱歉，创世区块已经存在，无法写入")
		return
	}
	coinbase,err:=transaction.NewCoinbaseTx(*address)
	if err != nil {
		fmt.Println("抱歉，构建coinbase交易遇到错误")
		return
	}
	client.Chain.Creatgenesis([]transaction.Transaction{*coinbase})
}
func (client *Client)SendTransaction(){
	fmt.Println("调用生成新区快的功能")
	addBlock:=flag.NewFlagSet(SENDTRASACtION,flag.ExitOnError)
	from:=addBlock.String("from","","发起者")
	to :=addBlock.String("to","","接收者")
	value:=addBlock.String("value","","数值")
	addBlock.Parse(os.Args[2:])
	fromSlice,err:=utils.JSONString2Slice(*from)
	toSlice,err:=utils.JSONString2Slice(*to)
	valueSlice,err:=utils.JSONFloat2Slice(*value)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	lenfrom:=len(fromSlice)
	lento:=len(toSlice)
	lenvalue:=len(valueSlice)
	if !(lenfrom ==lento && lenfrom==lenvalue) {
		fmt.Println("发起的交易参数不匹配")
		return
	}
	txs:=make([]transaction.Transaction,0)
	for index :=0;index<lenfrom ;index++  {
		tx,err:=transaction.NewTransaction(fromSlice[index],toSlice[index],valueSlice[index])
		if err != nil {
			fmt.Println("创建交易失败")
			return
		}
		txs=append(txs,*tx)
	}
	err=client.Chain.Addnewblock(txs)
	fmt.Println(fromSlice)
	fmt.Println(toSlice)
	fmt.Println(valueSlice)
	//tx,err:=transaction.NewTransaction(*from,*to,0)
	err =client.Chain.Addnewblock(txs)
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("已经成功发送交易")
}
func (client *Client)GetLastBlock(){
	set:=os.Args[2:]
	if len(set)>0 {
		fmt.Println("兄弟，你会错意了")
		return
	}
	fmt.Println("获取最新区块的功能")

	last:=client.Chain.GetLastBlock()
	hashBig:=new(big.Int)
	hashBig.SetBytes(last.Hash[:])
	if hashBig.Cmp(big.NewInt(0))>0 {
		fmt.Println("查询到最新区块")
		fmt.Println("最新区块高度：",last.Height)
		fmt.Println("最新区块的交易：",last.Txs)
		fmt.Printf("最新区块哈希%x\n",last.Hash)
		fmt.Printf("上一个区块哈希%x\n",last.PreHash)
		return
	}
	fmt.Println("抱歉，当前暂无最新区块")
	fmt.Println("请使用go run main.go generategensis生成创世区块")
}
func (client *Client)GetAllBlocks(){
	set:=os.Args[2:]
	if len(set)>0 {
		fmt.Println("抱歉,该功能不接收参数")
		return
	}
	fmt.Println("获取所有区块的功能")
	allBlocks,err:=client.Chain.GetAllblocks()
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	for _,block:=range allBlocks {
		fmt.Printf("区块%d，hash：%x，数据：%x\n",block.Height,block.Hash,block.Txs)
	}
}
func (client *Client)GetBlockCount(){
	fmt.Println("获取区块总数")
	blocks,err:=client.Chain.GetAllblocks()
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("查询成功，当前共有%d个区块\n",len(blocks))
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
	fmt.Println("\t"+SENDTRASACtION+"    创建一个新交易 可接收一个参数-data表示交易的数据")
	fmt.Println("\t"+GETALLBLOCKS+"      ")
	fmt.Println("\t"+GETLASTBLOCK+"   ")
	fmt.Println("\t"+GETBLOCKCOUNT+"  ")

	fmt.Println()
	fmt.Println("其他命令用法请使用help功能")
	fmt.Println()
}
