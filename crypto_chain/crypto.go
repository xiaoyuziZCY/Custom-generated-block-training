package crypto_chain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

//生成一对公私钥
func NewKey(curve elliptic.Curve)(*ecdsa.PrivateKey,error){
	curve =elliptic.P256()
	return ecdsa.GenerateKey(curve,rand.Reader)
}
func GetPub(curve elliptic.Curve,pri *ecdsa.PrivateKey)[]byte{
	return elliptic.Marshal(curve, pri.X, pri.Y)
}
