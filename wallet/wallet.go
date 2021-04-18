package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"github.com/boltdb/bolt"
)

const KEYSTORE="keystore"
//key
const ADDRESS ="address_keypair"
const COINBASE ="coinbase"
//钱包，用来管理地址和密钥对
type Wallet struct {
	Address map[string]*KeyPair
	Engine *bolt.DB
}

//定义创建新地址的方法
func (wallet *Wallet) CreateNewAddress() (string, error) {
	//1、获取秘钥对
	keyPair, err := NewKeyPair()
	if err != nil {
		return "", err
	}
	//2、将生成的keyPair中的公钥传递给NewAddress
	address, err := NewAddress(keyPair.Pub)
	if err != nil {
		return "", err
	}
	//放入到内存中
	wallet.Address[address] = keyPair
	//把新生成的地址和对应的KeyPair写入到DB中
	err = wallet.SaveMem2DB()
	return address, err
}

//保存内存中的地址和密钥对信息到DB文件里
func (wallet *Wallet)SaveMem2DB()error{
	var err error
	wallet.Engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(KEYSTORE))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(KEYSTORE))
			if err != nil {
				return err
			}
		}
		//for key,value:=range wallet.Address{
		//	keybytes:=bucket.Get([]byte(key))
		//	if len(keybytes)==0 {
		//		keyPairBytes,err:=utils.GobEncode(value)
		//		if err != nil {
		//			return err
		//		}
		//		bucket.Put([]byte(key),keyPairBytes)
		//	}
		//}
		gob.Register(elliptic.P256())
		buff:=new(bytes.Buffer)
		encoder:=gob.NewEncoder(buff)
		err:=encoder.Encode(wallet.Address)
		if err != nil {
			return err
		}
			bucket.Put([]byte(ADDRESS), buff.Bytes())
			return nil
	})
return err
}

//从db文件中加载数据，构建wallet结构体实例
func LoadWalletFromDB(db *bolt.DB)(Wallet,error){
	var wallet Wallet
	var err error
	adds :=make(map[string]*KeyPair)
db.View(func(tx *bolt.Tx) error {
	bucket:=tx.Bucket([]byte(KEYSTORE))
	if bucket==nil {
		return nil
	}
	addAndKeyPairBytes:=bucket.Get([]byte(ADDRESS))
	gob.Register(elliptic.P256())
	decoder:=gob.NewDecoder(bytes.NewReader(addAndKeyPairBytes))
	err:=decoder.Decode(&adds)
	if err != nil {
		return err
	}
	return nil
})
	//实例化wallet结构体，并给属性赋值
	wallet=Wallet{
		Address: adds,
		Engine:db,
	}
return wallet,err
}
//根据地址获取密钥对
func (wallet *Wallet)GetKeyPairByAddress(add string)*KeyPair{
	return wallet.Address[add]
}
func (wallet *Wallet) SetCoinbase(address string) error {
	//持久化的方式把用户设置的address进行保存
	var err error
	wallet.Engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(KEYSTORE))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(KEYSTORE))
			if err != nil {
				return err
			}
		}
		//keystore已经存在，可以将coinbase保存到桶keystore中
		bucket.Put([]byte(COINBASE), []byte(address))
		return nil
	})

	return err
}

/**
 * 该方法用于封装钱包的获取当前矿工地址的功能
 */
func (wallet *Wallet) GetCoinbase() string {
	var coinbase []byte

	wallet.Engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(KEYSTORE))
		if bucket == nil {
			return nil
		}
		//keystore桶存在，可以取coinbase地址
		coinbase = bucket.Get([]byte(COINBASE))
		return nil
	})
	return string(coinbase)
}
