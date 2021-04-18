package wallet

import (
	"Xianfeng/utils"
	"bytes"
)

const VERSION  ="0x00"
//生成一个新地址
func NewAddress(pub []byte)(string,error){
	//公钥推导地址
	hashpub:=utils.Hash256(pub)
	ripemdPub:=utils.RipEMd160(hashpub[:])
	versionPub:=append([]byte(VERSION),ripemdPub...)
	firstRipemd:=utils.Hash256(versionPub)
	secondRipemd:=utils.Hash256(firstRipemd[:])
	checkBytes:=secondRipemd[:4]
	origAddress :=append(versionPub,checkBytes...)
	address := utils.Encode(origAddress)
	return address,nil
}
//校验给定的一个字符串是否符合规范
func IsAddressValid(addr string)bool{
	//base58反编码
	reverseaddr:=utils.Decode(addr)
	if len(reverseaddr)<4{
		return false
	}
	//取后四字节作为校验位
	check:=reverseaddr[len(reverseaddr)-4:]
	//获取versionpub
	versionPub:=reverseaddr[:len(reverseaddr)-4]
	//双哈希
	firstHash:=utils.Hash256(versionPub)
	secondHash:=utils.Hash256(firstHash)
	reCheck:=secondHash[:4]
	//比较
	return bytes.Compare(check,reCheck)==0
}
