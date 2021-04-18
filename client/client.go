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
Chain chain.BlockChain
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
		client.CreateChain()
	case GENERATEGENESIS:
		client.GenerateGensis()
	case SENDTRANSACTION:
		client.SendTransaction()
	case GETBALANCE:
		client.GetBalance()
	case GETLASTBLOCK:
		client.GetLastBlock()
	case GETALLBLOCKS:
		client.GetAllBlocks()
	case GETNEWADDRESS:
		client.GetNewAddress()
	case GETBLOCKCOUNT:
		client.GetBlockCount()
	case LISTADDRESS://打印出当前的地址列表
		client.ListAddress()
	case SETCOINBASE:
		client.SetCoinbase()
	case DUMPPRIVATEKEY:
		client.DumpprivateKey()
	case HELP:
		fmt.Println("获取使用说明")
	default:
		fmt.Println("抱歉，目前不支持该命令功能")
		fmt.Println("请你使用help功能参考说明")

	}
}
func (client *Client)GenerateGensis(){
	fmt.Println("调用生成创世区块的功能")
	generateGensis:=flag.NewFlagSet(GENERATEGENESIS,flag.ExitOnError)
	address:=generateGensis.String("address","","用户指定的矿工地址")
	generateGensis.Parse(os.Args[2:])
	//判断是否已经存在创世区块
	hashBig:=new(big.Int)
	hashBig.SetBytes(client.Chain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0))==1 {
		fmt.Println("抱歉，已有coinbase交易")
		return
	}
	coinbaseHash,err:=client.Chain.CreateCoinbase(*address)
	if err != nil {
		fmt.Println("抱歉，coinbase交易构建失败",err.Error())
		return
	}
	fmt.Printf("恭喜，COINBASE交易创建成功，hash：%x\n",coinbaseHash)
}
func (client *Client)SendTransaction(){
	fmt.Println("调用生成新交易的功能")
	addBlock:=flag.NewFlagSet(SENDTRANSACTION,flag.ExitOnError)
	from:=addBlock.String("from","","发起者")
	to :=addBlock.String("to","","接收者")
	value:=addBlock.String("value","","数值")
	//label :=addBlock.String("label","","转账备注")
	addBlock.Parse(os.Args[2:])
	err := client.Chain.SendTransaction(*from, *to, *value)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
fmt.Println("成功发送交易")

}
//获取地址的余额
func (client * Client)GetBalance(){
	var address string
	getbalance:=flag.NewFlagSet(GETBALANCE,flag.ExitOnError)
	getbalance.StringVar(&address,"address","","要查询的地址")
	getbalance.Parse(os.Args[2:])
	if len(address)==0 {
		fmt.Println("抱歉，请输入要查询的地址")
		return
	}
	totalbalance:=client.Chain.GetBalance(address)
	fmt.Printf("地址%s余额为：%f\n",address,totalbalance)
}
func (client *Client) SetCoinbase() {
	setCoinbase := flag.NewFlagSet(SETCOINBASE, flag.ExitOnError)
	coinbase := setCoinbase.String("address", "", "用户设置的矿工地址")
	setCoinbase.Parse(os.Args[2:])
	if len(os.Args[2:]) > 2 {
		fmt.Println("抱歉，不支持的命令行参数，请重试！")
		return
	}
	client.Chain.SetCoinbase(*coinbase)
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
	allBlocks,err:=client.Chain.GetAllBlocks()
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	for _,block:=range allBlocks {
		fmt.Printf("区块%d，hash：%x，数据：%x\n",block.Height,block.Hash,block.Txs)
		for index,tx:=range block.Txs  {
			fmt.Printf("区块%d的第%d笔交易,交易hash是：%x\n", block.Height, index, tx.TxHash)
			fmt.Println("\t该笔交易的交易输入：")
			for inputIndex, input := range tx.Inputs { //遍历交易的交易输入
				fmt.Printf("\t\t第%d个交易输入,花的钱来自%x中的第%d个输出\n", inputIndex, input.TxId, input.Vout)
			}
			fmt.Println("\t该笔交易的交易输出：")
			for outputIndex, output := range tx.Outputs { //遍历交易的交易输出
				fmt.Printf("\t\t第%d个交易输出，转给%x一笔面额为%f的钱\n", outputIndex, output.PubKHash, output.Value)
			}
		}
		fmt.Println()
	}
}
//导出某个地址的私钥
func (client *Client)DumpprivateKey(){
	dumpPrivateKey:=flag.NewFlagSet(DUMPPRIVATEKEY,flag.ExitOnError)
	address:=dumpPrivateKey.String("address","","")
	dumpPrivateKey.Parse(os.Args[2:])
	if len(os.Args[2:])>2 {
		fmt.Println("抱歉，您输入的参数有误")
		return
	}
	privateKey,err:=client.Chain.DumpPrivateKey(*address)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("私钥是：",privateKey.D.Bytes())
}
//用于列出当前节点的地址
func (client *Client)ListAddress(){
	listAddress:=flag.NewFlagSet(LISTADDRESS,flag.ExitOnError)
	listAddress.Parse(os.Args[2:])
	if len(os.Args[2:])>0 {
		fmt.Println("抱歉，无法解析打印地址列表参数")
		return
	}
	addList,err:=client.Chain.GetAddressList()
	if err != nil {
		fmt.Println("加载地址遇到错误：",err.Error())
		return
	}
	if len(addList)==0 {
		fmt.Println("当前无节点")
		return
	}
	fmt.Println("地址列表获取成功")
	for index,address:=range addList {
		fmt.Printf("(%d):%s\n",index+1,address)
	}
}
//生成一个新地址
func (Client *Client)GetNewAddress(){
	getNewaddress:=flag.NewFlagSet(GETNEWADDRESS,flag.ExitOnError)
	err:=getNewaddress.Parse(os.Args[2:])
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	//判断用户是否输了内容
	if len(os.Args[2:])>0 {
		fmt.Println("生成地址不需要参数，请重试")
		return
	}
	address,err:=Client.Chain.GetNewAddress()
	if err != nil {
		fmt.Println("地址生成错误：",err.Error())
		return
	}
	fmt.Println("生成的地址为：",address)
}
func (client *Client)GetBlockCount(){
	fmt.Println("获取区块总数")
	blocks,err:=client.Chain.GetAllBlocks()
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("查询成功，当前共有%d个区块\n",len(blocks))
}

func (client *Client)CreateChain() {

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
	fmt.Println("\t"+SENDTRANSACTION+"    创建一个新交易 可接收一个参数-data表示交易的数据")
	fmt.Println("\t"+GETALLBLOCKS+"      ")
	fmt.Println("\t"+GETLASTBLOCK+"   ")
	fmt.Println("\t"+GETBLOCKCOUNT+"  ")
	fmt.Println("\t"+GETNEWADDRESS+"     ")
	fmt.Println("\t"+LISTADDRESS+"       ")
	fmt.Println()
	fmt.Println("其他命令用法请使用help功能")
	fmt.Println()
}
