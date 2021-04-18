package consensus

import (
	"fmt"
)

type POS struct {
Block BlockInterface
}
func (pos POS) SearchNonce()([32]byte,int64){
	fmt.Println("已为pos算法机制")
	return [32]byte{},0
}
