package loan

import (
	"fmt"
	"testing"
)

func TestEqualInterest(t *testing.T) {
	EqualInterest(2000000, 30)
	EqualInterest(2000000, 20)
	EqualInterest(3000000, 30)
}

func TestCalcMonth(t *testing.T) {
	CalcMonth(3000000, 1000000, 3)
	fmt.Println()
	CalcMonth(2000000, 1000000, 2)
}
