package consensus

import (
	"Xianfeng/utils"
	"bytes"
	"crypto/sha256"
	"math/big"
)

// 256位二进制
//思路：给一个大整数，初始值为1，根据自己需要的难度进行左移位，左移的位数是256-0的个数

//000000010000000....000

const DIFFICULTY = 16 //初始难度为10，即大整数的开头有10个零

/**
 * 工作量证明
 */
type ProofWork struct {
	Block  BlockInterface
	Target *big.Int
}

/**
 * 实现共识机制接口的方法
 */
func (work ProofWork) SearchNonce() ([32]byte, int64) {
	//block -> nonce
	// block哈希 < 系统提供的某个目标值
	//1 给定一个non值，计算带有non的区块哈希
	var nonce int64
	nonce = 0
	hashBig := new(big.Int)
	for {
		hash := CalculateBlockHash(work.Block, nonce)
		//2 系统给定的值
		target := work.Target
		//3 拿1和2比较
		//hash [32]byte
		//target big.Int

		//result := bytes.Compare(hash[:], target.Bytes())
		hashBig = hashBig.SetBytes(hash[:])
		result := hashBig.Cmp(target)
		//4 判断结果，区块哈希<给定值，返回non;
		if result == -1 {
			return hash, nonce
		}
		//否则，non自增
		nonce++
	}
}

/**
 * 根据当前的区块和当前的non值，计算区块的哈希值
 */
func CalculateBlockHash(block BlockInterface, nonce int64) [32]byte {
	heightByte, _ := utils.Int2Byte(block.GetHeight())
	versionByte, _ := utils.Int2Byte(block.GetVersion())
	timeByte, _ := utils.Int2Byte(block.GetTimeStamp())
	nonceByte, _ := utils.Int2Byte(nonce)
	preHash := block.GetPreHash()
	txs := block.GetTxs()
	//Transaction -> []byte 序列化
	//gob -> 区块的序列化
	txsBytes := make([]byte, 0)
	for _, tx := range txs {
		txBytes, err := utils.GobEncode(tx)
		if err != nil {
			break
		}
		//buff.Bytes()//每一个交易的序列化数据
		txsBytes = append(txsBytes, txBytes...)
	}

	//bytes.Join
	bk := bytes.Join([][]byte{heightByte,
		versionByte,
		preHash[:],
		timeByte,
		nonceByte,
		txsBytes,
	},
		[]byte{})
	return sha256.Sum256(bk)
}