package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"golang.org/x/crypto/ripemd160"
)

//将int数值类型转换为byte
func Int2Byte(num int64)([]byte,error){
	buff :=new(bytes.Buffer)
	err :=binary.Write(buff,binary.BigEndian,num)
	if err !=nil {
		return nil,err
	}
	return buff.Bytes(),nil
}
//公共的gob序列化
func GobEncode(entity interface{})([]byte,error){
	buff:=new(bytes.Buffer)
	encoder:=gob.NewEncoder(buff)
	err:=encoder.Encode(entity)
	return buff.Bytes(),err
}
//公共的gob反序列化
func GobDecode(data []byte,entity interface{})(interface{},error){
	decoder:=gob.NewDecoder(bytes.NewReader(data))
	err:=decoder.Decode(entity)
	return entity,err
}
func JSONString2Slice(data string)([]string,error)  {
	var slice []string
	err :=json.Unmarshal([]byte(data),&slice)
	return slice,err
}
func JSONFloat2Slice(data string)([]float64,error)  {
	var slice []float64
	err :=json.Unmarshal([]byte(data),&slice)
	return slice,err
}
func Hash256(data []byte)([]byte){
	sha256Hash := sha256.New()
	sha256Hash.Write(data)
	return sha256Hash.Sum(nil)
}
func RipEMd160(data []byte)([]byte){
	ripemdHash:=ripemd160.New()
	ripemdHash.Write(data)
	return ripemdHash.Sum(nil)
}