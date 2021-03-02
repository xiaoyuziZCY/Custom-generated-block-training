package pos

import "fmt"

type POS struct {

}
func (pos POS) Run() interface{}{
	fmt.Println("已为pos算法机制")
	return nil
}
