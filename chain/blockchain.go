package chain

import (
	"Xianfeng/transaction"
	"Xianfeng/utils"
	"Xianfeng/wallet"
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"math/big"
)

const BLOCKS = "blocks"
const LASTHASH = "lastHash"

/**
 * 定义区块链这个结构体，用于存储产生的区块（内存中)
 */
type BlockChain struct {
	//Blocks []Block
	Engine            *bolt.DB
	LastBlock         Block          //最新的区块
	IteratorBlockHash [32]byte       //迭代到的区块哈希值
	Wallet            *wallet.Wallet //钱包管理模块中的Wallet
}

func NewBlockChain(db *bolt.DB) (BlockChain, error) {
	//增加为lastblock赋值的逻辑
	var lastBlock Block
	var err error
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(BLOCKS))
			if err != nil {
				return err
			}
		}
		lastHash := bucket.Get([]byte(LASTHASH))
		if len(lastHash) == 0 {
			return nil
		}
		lastBlockBytes := bucket.Get(lastHash)
		lastBlock, err = Deserialize(lastBlockBytes)
		if err != nil {
			return err
		}
		return nil
	})
	blockChain := BlockChain{
		Engine:            db,
		LastBlock:         lastBlock,
		IteratorBlockHash: lastBlock.Hash,
	}
	wlt, err := wallet.LoadWalletFromDB(db)
	if err != nil {
		return blockChain, err
	}
	//把构建的wallet对象赋值给blockChain的Wallet属性变量
	blockChain.Wallet = &wlt
	return blockChain, err
}

/**
 * 创建一个区块链实例，该实例携带一个创世区块
 */
func (chain *BlockChain) CreateGenesis(txs []transaction.Transaction) {
	//先看chain.LastBlock是否为空
	hashBig := new(big.Int)
	hashBig.SetBytes(chain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) > 0 {
		return
	}

	engine := chain.Engine
	//读一遍bucket，查看是否有数据
	engine.Update(func(tx *bolt.Tx) error { //
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			bucket, _ = tx.CreateBucket([]byte(BLOCKS))
		}
		if bucket != nil {
			lastHash := bucket.Get([]byte(LASTHASH))
			if len(lastHash) == 0 {
				genesis := CreatGenesisBlock(txs)
				genSerBytes, _ := genesis.Serialize()
				//存创世区块
				bucket.Put(genesis.Hash[:], genSerBytes)
				//更新最新区块的标志 lastHash -> 最新区块hash
				bucket.Put([]byte(LASTHASH), genesis.Hash[:])
				chain.LastBlock = genesis
				chain.IteratorBlockHash = genesis.Hash
			}
		}
		return nil
	})
}

/**
 * 该方法用于创建一笔coinbase交易
 */
func (chain *BlockChain) CreateCoinbase(addr string) ([]byte, error) {
	//1、判断有效性
	isValid := wallet.IsAddressValid(addr)
	if !isValid {
		return nil, errors.New("地址不合法，请检查后重试！")
	}
	//2、创建coinbase交易
	coinbase, err := transaction.NewCoinbaseTx(addr)
	if err != nil {
		return nil, err
	}
	chain.CreateGenesis([]transaction.Transaction{*coinbase})
	return coinbase.TxHash[:], nil
}

/**
 * 获取某个地址的余额
 */
func (chain *BlockChain) GetBalance(addr string) float64 {
	_, totalbalance := chain.GetUTXOsWithBalance(addr, []transaction.Transaction{})
	return totalbalance
}

/**
 * 获取某个特定地址的余额和所能花费的utxo集合
 */
func (chain *BlockChain) GetUTXOsWithBalance(addr string, txs []transaction.Transaction) ([]transaction.UTXO, float64) {
	//1、从文件中遍历区块，找出区块中已经存在交易中的可花费utxo
	dbUtxos := chain.SearchUTXOs(addr)
	fmt.Println(dbUtxos)
	//2、遍历内存中的txs切片, 如果当前已构建还未存储的交易已经花了前，要剔除掉
	memSpends := make([]transaction.TxInput, 0)
	memInComes := make([]transaction.UTXO, 0)
	keyPair := chain.Wallet.GetKeyPairByAddress(addr)
	if keyPair == nil {
		return nil, 0
	}

	for _, tx := range txs {
		//花的钱
		for _, input := range tx.Inputs {
			if input.CheckPubKWithInput(keyPair.Pub) {
				memSpends = append(memSpends, input)
			}
		}
		//收入的钱
		for index, output := range tx.Outputs {
			if output.VertifyOutputWithAddress(addr) {
				utxo := transaction.NewUTXO(tx.TxHash, index, output)
				memInComes = append(memInComes, utxo)
			}
		}
	}

	//3、合并1和2, 将内存中已经花掉的utxo从dbUtxo删除掉，将内存中产生的收入加入到可花费收入中
	utxos := make([]transaction.UTXO, 0)
	var isSpend bool
	for _, dbUtxo := range dbUtxos {
		isSpend = false
		for _, memInput := range memSpends {
			//判断某一个utxo是否已经被消费了
			if dbUtxo.IsSpent(memInput) {
				isSpend = true
			}
		}
		if !isSpend {
			utxos = append(utxos, dbUtxo)
		}
	}
	//把内存中的产生的收入放入到可花的utxo中
	utxos = append(utxos, memInComes...)

	var totalBalance float64
	for _, utxo := range utxos {
		totalBalance += utxo.Value
	}
	return utxos, totalBalance
}

/**
 * 发送交易的功能方法
 */
func (chain *BlockChain) SendTransaction(from string, to string, value string) (error) {
	fromSlice, err := utils.JSONString2Slice(from)
	toSlice, err := utils.JSONString2Slice(to)
	valueSlice, err := utils.JSONFloat2Slice(value)
	if err != nil {
		return err
	}

	//判断参数的长度，筛选参数不匹配的情况
	lenFrom := len(fromSlice)
	lenTo := len(toSlice)
	lenValue := len(valueSlice)
	if !(lenFrom == lenTo && lenFrom == lenValue) {
		return errors.New("发起交易的参数不匹配，请检查后重试")
	}

	//地址有效性的判断
	for i := 0; i < len(fromSlice); i++ {
		//交易发起人的地址是否合法，合法为true，不合法为false
		isFromValid := wallet.IsAddressValid(fromSlice[i])
		//交易接收者的地址是否合法，合法为true，不合法为false
		isToValid := wallet.IsAddressValid(toSlice[i])
		//from: 合法   合法
		//to:   不合法  不合法
		if !isFromValid || !isToValid {
			return errors.New("交易的参数地址不合法，请检查后重试")
		}
	}

	//遍历参数的切片，创建交易
	txs := make([]transaction.Transaction, 0)
	for index := 0; index < lenFrom; index++ {
		utxos, totalBalance := chain.GetUTXOsWithBalance(fromSlice[index], txs)
		fmt.Printf("转账发起人%s,当前余额：%f,接收者:%s,转账数额：%f\n", fromSlice[index], totalBalance, toSlice[index], valueSlice[index])
		if totalBalance < valueSlice[index] {
			return errors.New("抱歉，" + fromSlice[index] + "余额不足，请充值！")
		}
		var inputAmount float64 //总的花费的钱数

		utxoNum := 0
		for num, utxo := range utxos {
			inputAmount += utxo.Value
			if inputAmount >= valueSlice[index] {
				//够花了
				utxoNum = num
				break
			}
		}
		keyPair := chain.Wallet.GetKeyPairByAddress(fromSlice[index])
		//1、创建交易
		tx, err := transaction.NewTransaction(
			utxos[:utxoNum+1],
			fromSlice[index],
			keyPair.Pub,
			toSlice[index],
			valueSlice[index])
		if err != nil {
			return errors.New("抱歉，创建交易失败，请检查后重试")
		}
		//2、使用from对应的私钥对tx进行交易签名
		//keyPair := chain.Wallet.GetKeyPairByAddress(fromSlice[index])

		err =tx.Sign(keyPair.Priv, utxos[:utxoNum+1])
		if err != nil {
			return err
		}

		txs = append(txs, *tx)
	}
	address:=chain.GetCoinbase()
	if len(address)==0{
		return errors.New("暂未设置矿工地址")
	}
	coinbase,err:=transaction.NewCoinbaseTx(address)
	if err!=nil {
		return err
	}
	sumTxs:=make([]transaction.Transaction,0)
	sumTxs=append(sumTxs,*coinbase)
	sumTxs=append(sumTxs,txs...)
	for _,tx:=range txs {
		if tx.IsCoinbaseTransaction() {
			continue
		}
		spendUTXOs,err:=chain.FindSpentUTXOsByTx(tx,txs)
		if err != nil {
			return err
		}
		verify,err:=tx.VerifySign(spendUTXOs)
		if err != nil {
			return err
		}
		if !verify {
			return errors.New("交易失败请重试")
		}
	}
	//把构建好的交易存入到区块中
	err = chain.AddNewBlock(txs)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 新增一个区块
 */
func (chain *BlockChain) AddNewBlock(txs []transaction.Transaction) error {
	//1、从db中找到最后一个区块数据
	engine := chain.Engine
	//2、获取到最新的区块
	lastBlock := chain.LastBlock

	//3、得到最后一个区块的各种属性，并利用这些属性生成新区块
	newBlock := CreateBlock(lastBlock.Height, lastBlock.Hash, txs)
	newBlockByte, err := newBlock.Serialize()
	if err != nil {
		return err
	}
	//4、更新db文件，将新生成的区块写入到db中，同时更新最新区块的指向标记
	engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			err = errors.New("区块链数据库操作失败，请重试!")
			return err
		}
		//将最新的区块数据存到db中
		bucket.Put(newBlock.Hash[:], newBlockByte)
		//更新最新区块的指向标记
		bucket.Put([]byte(LASTHASH), newBlock.Hash[:])

		//更新blockChain对象的LastBlock结构体实例
		chain.LastBlock = newBlock
		chain.IteratorBlockHash = newBlock.Hash
		return nil
	})
	return err
}

/**
 * 获取最新的最后的一个区块
 */
func (chain *BlockChain) GetLastBlock() (Block) {
	return chain.LastBlock
}

/**
 * 获取所有的区块
 */
func (chain BlockChain) GetAllBlocks() ([]Block, error) {
	engine := chain.Engine
	var errs error
	blocks := make([]Block, 0)
	engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			errs = errors.New("区块数据库操作是吧，请重试！")
			return errs
		}

		var currentHash []byte
		//直接从倒数第二个区块进行遍历
		currentHash = bucket.Get([]byte(LASTHASH))
		for { //倒数第一个区块区块开始遍历
			//根据区块hash拿[]byte类型的区块数据
			currentBlockBytes := bucket.Get(currentHash)
			//[]byte类型的区块数据 反序列化为  struct类型
			currentBlock, err := Deserialize(currentBlockBytes)
			if err != nil {
				errs = err
				break
			}
			blocks = append(blocks, currentBlock)
			//终止循环的逻辑
			if currentBlock.Height == 0 {
				break
			}
			//创世区块的hash值
			currentHash = currentBlock.PreHash[:]
		}
		return nil
	})
	return blocks, errs
}

/**
 * 该方法用于实现ChainIterator迭代器接口的方法，用于判断是否还有区块
 */
func (chain *BlockChain) HasNext() bool {
	//是否还有前一个区块
	//思路：先知道当前在哪个区块，根据当前的区块去判断是否还有下一个区块
	engine := chain.Engine
	var hasNext bool
	engine.View(func(tx *bolt.Tx) error {
		currentBlockHash := chain.IteratorBlockHash
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块数据文件操作失败,请重试")
		}
		currentBlockBytes := bucket.Get(currentBlockHash[:])
		currentBlock, err := Deserialize(currentBlockBytes)
		if err != nil {
			return err
		}

		hashBig := big.NewInt(0)
		hashBig = hashBig.SetBytes(currentBlock.Hash[:])
		if hashBig.Cmp(big.NewInt(0)) == 1 { //区块hash有值
			hasNext = true
		} else {
			hasNext = false
		}

		return nil
	})
	return hasNext
}

/**
 * 该方法用于实现ChainIterator迭代器接口的方法，用于取出下一个区块
 */
func (chain *BlockChain) Next() Block {
	engine := chain.Engine
	var currentBlock Block
	engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块数据文件操作失败,请重试！")
		}
		currentBlockBytes := bucket.Get(chain.IteratorBlockHash[:])
		currentBlock, _ = Deserialize(currentBlockBytes)
		chain.IteratorBlockHash = currentBlock.PreHash //赋值iteratorBlock，
		return nil
	})
	return currentBlock
}

/**
 * 定义该方法，用于实现寻找与from有关的所有可花费的交易输出，即寻找UTXO
 */
func (chain BlockChain) SearchUTXOs(from string) ([]transaction.UTXO) {

	//定义容器，存放from的所有的花费
	spends := make([]transaction.TxInput, 0)

	//定义容器，存放from的所有的收入
	inComes := make([]transaction.UTXO, 0)

	keyPair := chain.Wallet.GetKeyPairByAddress(from)
	if keyPair == nil {
		return nil
	}
	//使用迭代器进行区块的遍历
	for chain.HasNext() { //遍历区块
		block := chain.Next()
		for _, tx := range block.Txs { //遍历区块的交易
			//a、遍历交易输入
			for _, input := range tx.Inputs {
				if !input.CheckPubKWithInput(keyPair.Pub) {
					continue
				}
				//该交易输入是from的，即from花钱了
				spends = append(spends, input)
			}
			//b、遍历交易输出
			for index, output := range tx.Outputs {
				if !output.VertifyOutputWithAddress(from) {
					continue
				}
				//该交易输出是from的，即from有收入
				input := transaction.NewUTXO(tx.TxHash, index, output)
				inComes = append(inComes, input)
			}
		}
	}

	utxos := make([]transaction.UTXO, 0)
	//遍历spends和inComes,将已花费的记录剔除掉，剩下可花费的UTXO
	var isInComeSpend bool
	for _, income := range inComes {
		//判断每一笔收入是否在之前的交易中已经被花过了
		isInComeSpend = false
		for _, spend := range spends { //5
			if income.Txid == spend.TxId && income.Vout == spend.Vout {
				isInComeSpend = true
				break
			}
		}
		//追加
		if !isInComeSpend { //isInComeSpend如果如果为false，表示未被花,可加到utxos中
			utxos = append(utxos, income)
		}
	}
	return utxos
}

/**
 * 生成一个新的地址，并返回
 */
func (chain *BlockChain) GetNewAddress() (string, error) {
	chain.Wallet.Engine = chain.Engine
	return chain.Wallet.CreateNewAddress()
}

/**
 * 该方法用于获取到当前节点管理的地址，并返回
 */
func (chain *BlockChain) GetAddressList() ([]string, error) {
	if chain.Wallet == nil {
		return nil, errors.New("钱包报错，请重试")
	}
	addList := make([]string, 0)
	for address, _ := range chain.Wallet.Address {
		addList = append(addList, address)
	}
	return addList, nil
}

/**
 * 该方法用于导出某个特定地址的私钥
 */
func (chain *BlockChain) DumpPrivateKey(add string) (*ecdsa.PrivateKey, error) {
	isAddrValid := wallet.IsAddressValid(add)
	if !isAddrValid {
		return nil, errors.New("地址不合法，请检查后重试")
	}

	if chain.Wallet == nil {
		return nil, errors.New("钱包报错，请重试！")
	}

	//Address 是一个map，可以使用map[key]
	keyPair := chain.Wallet.Address[add]
	if keyPair == nil {
		return nil, errors.New("抱歉,未找到该地址的私钥")
	}
	return keyPair.Priv, nil
}
//用于接收特点的交易头像，寻找到某个交易具体花了哪些utxo
func (chain BlockChain)FindSpentUTXOsByTx(transac transaction.Transaction,memtxs []transaction.Transaction)([]transaction.UTXO,error)  {
	var err error
	spentUTXOs:=make([]transaction.UTXO,0)
	for chain.HasNext() {
		block:=chain.Next()
		chain.Engine.View(func(tx *bolt.Tx) error {
			bucket:=tx.Bucket([]byte(BLOCKS))
			if bucket==nil {
				err=errors.New("未找到交易")
				return err
			}
			for _,tran:=range block.Txs {
				//遍历input
				for outIndex,output:=range tran.Outputs{
					utxo:=transaction.NewUTXO(tran.TxHash,outIndex,output)
					for _,input:=range transac.Inputs  {
						if utxo.IsSpent(input){
							spentUTXOs=append(spentUTXOs,utxo)
						}
					}
				}
			}
			return nil
		})
	}
	//遍历内存中的交易容器，找出此笔交易花费的utxo
	for _,memtx:=range memtxs {
		if bytes.Compare(memtx.TxHash[:],transac.TxHash[:])==0 {
			continue
		}
		for outIndex,output:=range memtx.Outputs {
			utxo:=transaction.NewUTXO(memtx.TxHash,outIndex,output)
			for _,input:=range transac.Inputs  {
				if utxo.IsSpent(input) {
					spentUTXOs=append(spentUTXOs,utxo)
				}
			}
		}
	}
	return spentUTXOs,nil
}
func (chain *BlockChain) SetCoinbase(address string) error {
	//1、先做地址的规范性校验
	if !wallet.IsAddressValid(address) {
		return errors.New("检查地址的规范性，请重试！")
	}
	return chain.Wallet.SetCoinbase(address)
}


func (chain *BlockChain) GetCoinbase() string {
	return chain.Wallet.GetCoinbase()
}