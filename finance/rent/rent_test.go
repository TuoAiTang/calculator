package rent

import (
	"fmt"
	"testing"
)

func Test_rent(t *testing.T) {
	a := rent(6130, 0)
	b := rent(6530, 0)
	c := rent(7490, 1500)
	fmt.Printf("租金(6130):%.2f,对比6530:%.2f,对比7490:%.2f\n", a, b-a, c-a)
	fmt.Printf("租金(6530):%.2f,对比7490:%.2f\n", b, c-b)
	fmt.Printf("租金(7490，优惠1500):%.2f\n", c)
}
