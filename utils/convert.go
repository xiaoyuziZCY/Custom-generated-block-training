package utils

import (
	"bytes"
	"encoding/binary"
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
