package pension

import (
	"fmt"
	"strings"
	"testing"
)

func TestCalc(t *testing.T) {
	p := Param{
		RetireCitySalaryAvg:   ChengduAvgSalaryYear / 12,
		PayYear:               15,
		SalaryAvgBeforeRetire: 6500,
		ExpectAge:             60,
		LifeYearAvg:           AvgLifeYear,
		PayUpLimit:            BeijingUpLimit,
	}
	r := Calc(p)
	fmt.Printf("参数:\n%s\n%s\n", p, strings.Repeat("=", 50))
	fmt.Printf("结果:\n%s\n%s\n", r, strings.Repeat("=", 50))
}
