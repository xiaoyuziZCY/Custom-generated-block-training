package pow

import "fmt"

type POW struct {

}
func (pow POW) Run() interface{}{
	fmt.Println("已为pow算法机制")
	return nil
}
