package rent

import (
	"fmt"
	"testing"
)

func Test_rent(t *testing.T) {
	fmt.Printf("租金(6130):%.2f\n", rent(6130, 0))
	fmt.Printf("租金(6630):%.2f\n", rent(6630, 0))
	fmt.Printf("租金(7490，优惠1500):%.2f\n", rent(7490, 1500))
}
