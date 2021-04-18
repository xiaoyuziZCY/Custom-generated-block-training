package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

//定义wallet一个钱包
type KeyPair struct {
	Priv *ecdsa.PrivateKey
	Pub []byte
}


func NewKeyPair()(*KeyPair,error){
	curve:=elliptic.P256()
	pri,err:=ecdsa.GenerateKey(curve,rand.Reader)
	if err!=nil {
		return nil,err
	}
	pub:=elliptic.Marshal(curve,pri.X,pri.Y)
	keyPair:=&KeyPair{
		Priv: pri,
		Pub:  pub,
	}
	return keyPair,nil
}
//[]byte类型数据转pubklickey类型公钥
func GetPublicKeyWithBytes(curve elliptic.Curve, data []byte) ecdsa.PublicKey {
	x, y := elliptic.Unmarshal(curve, data)
	return ecdsa.PublicKey{curve, x, y}
}

/**
 * 将[]byte类型签名恢复为内存中的*big.Int类型
 */
func RestoreSignature(sign []byte) (r, s *big.Int) {
	//rBytes := sign[:len(sign)/2]
	//sBytes := sign[len(sign)/2:]

	rBig := new(big.Int)
	sBig := new(big.Int)

	rBig.SetBytes(sign[:len(sign)/2])
	sBig.SetBytes(sign[len(sign)/2:])

	return rBig, sBig
}
