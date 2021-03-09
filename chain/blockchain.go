package chain

import (
	"bytes"
	"errors"
	"github.com/boltdb/bolt"
	"math/big"
)

const BLOCKS  ="blocks"
const LASTHASH  = "lasthash"
//用于存储在内存里产生的区块
type Blockchain struct {
	//文件操作对象
	Engine *bolt.DB
	LastBlock Block//最新区块
	IteratorBlockHash [32]byte//迭代到的区块hash
}
func Newblockchain(db *bolt.DB)Blockchain{
	return Blockchain{Engine:db}
}
//创建一个区块链，携带创世区块
func (chain *Blockchain)Creatgenesis(genesisdata []byte) {
	engine := chain.Engine
	//读一遍是否有数据
	engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil { //为空证明是第一次
			bucket, _ = tx.CreateBucket([]byte(BLOCKS))
		}
		if bucket != nil {
			lasthash := bucket.Get([]byte(LASTHASH))
			if len(lasthash) == 0 {
				genesis := CreatGenesisBlock(genesisdata)
				genserbytes, _ := genesis.Serialize()
				//存创世区块
				bucket.Put(genesis.Hash[:], genserbytes)
				//更新最新区块的标志
				bucket.Put([]byte(LASTHASH), genesis.Hash[:])
				chain.LastBlock=genesis
				chain.IteratorBlockHash=genesis.Hash
			} else {
				//创世区块不存在了,不需要写入了
				lasthash := bucket.Get([]byte(LASTHASH))
				lastBlockbytes := bucket.Get(lasthash)
				lastBlock, _ := Deserialize(lastBlockbytes)
				chain.LastBlock = lastBlock
				chain.IteratorBlockHash=lastBlock.Hash
			}
		}
		return nil
	})
}

func (chain *Blockchain) Addnewblock(data []byte)error{
	//从db找最后区块
	engine :=chain.Engine
	lastBlock:=chain.LastBlock
			newblock:=CreateBlock(lastBlock.Height,lastBlock.Hash,data)
			newblockByte,err:=newblock.Serialize()
			if err !=nil {
				return err
			}
		engine.Update(func(tx *bolt.Tx) error {
			bucket:=tx.Bucket([]byte(BLOCKS))
			if bucket==nil {
				err =errors.New("区块链数据操作失败")
				return err
			}
			bucket.Put(newblock.Hash[:],newblockByte)
			bucket.Put([]byte(LASTHASH),newblock.Hash[:])
			//
			chain.LastBlock=newblock
			chain.IteratorBlockHash=newblock.Hash
			return nil
		})
		return err
		}
		//获取最新一个区块
func (chain Blockchain)GetLastBlock()(Block){
	return chain.LastBlock
}
//获取所有的区块
func(chain Blockchain)GetAllblocks()([]Block,error){
	engine:=chain.Engine
	var errs error
	genesishash := [32]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	blocks:=make([]Block,0)
	errs = engine.View(func(tx *bolt.Tx) error {
		bucket:=tx.Bucket([]byte(BLOCKS))
		if bucket==nil {
			errs =errors.New("块数据库操作失败")
			return errs
		}
		////将最后的区块存储到切片里
		blocks=append(blocks,chain.LastBlock)
		var currentHash []byte
		//直接从倒数第二个区块进行遍历
		currentHash = chain.LastBlock.PreHash[:]
		for {
			//根据区块hash拿[]byte类型的区块数据
			currentBlockBytes:=bucket.Get(currentHash)
			//[]byte类型区块数据反序列化
			currentBlock,err:=Deserialize(currentBlockBytes)
			if err !=nil {
				errs=err
				break
			}
			blocks=append(blocks,currentBlock)

			if bytes.Compare(currentHash, genesishash[:]) == 0 {
				break
			}
			currentHash=currentBlock.PreHash[:]
		}
		return nil
	})
	return blocks,errs
}
//用以实现ChainIterator迭代器接口方法，判断是否还有区块
func (chain *Blockchain)HasNext()bool{
	//先知道当前在哪个区块，然后判断是否有下一个区块
	engine :=chain.Engine
	var  hashnext bool
	engine.View(func(tx *bolt.Tx) error {
		currentBlockHash:=chain.IteratorBlockHash
		bucket :=tx.Bucket([]byte(BLOCKS))
		if bucket==nil {
			return errors.New("区块数据文件操作失败")
		}
		currentBlockBytes:=bucket.Get(currentBlockHash[:])
		currentBlock,err:=Deserialize(currentBlockBytes)
		if err!=nil {
			return err
		}
		hashBig:=big.NewInt(0)
		hashBig=hashBig.SetBytes(currentBlock.Hash[:])
		if hashBig.Cmp(big.NewInt(0))>0 {
			hashnext=true
		}else {
			hashnext=false
		}
		//preBlockBytes:=bucket.Get(currentBlock.PreHash[:])
		//hasnext=len(preBlockBytes)!= 0
		return nil
	})
	return hashnext
}
//用以实现ChainIterator迭代器接口方法，取出下一个区块
func (chain *Blockchain)Next()Block{
	engine:=chain.Engine
	var currentBlock Block
	engine.View(func(tx *bolt.Tx) error {
		bucket:=tx.Bucket([]byte(BLOCKS))
		if bucket==nil {
			return errors.New("区块数据文件操作失败")
		}
		currentBlockBytes:=bucket.Get(chain.IteratorBlockHash[:])
		currentBlock,_=Deserialize(currentBlockBytes)
		chain.IteratorBlockHash=currentBlock.PreHash//赋值
		return nil
	})
	return currentBlock
}