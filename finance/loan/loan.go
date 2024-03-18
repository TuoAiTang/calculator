package loan

import (
	"fmt"
	"math"
)

// 等额本息: 每月还款额=[贷款本金×月利率×（1+月利率）^还款月数]÷[（1+月利率）^还款月数－1]
func EqualInterest(principal float64, year float64) (monthRepay float64) {
	month := year * 12
	monthRate := 0.04 / 12
	monthRepay = principal * monthRate * math.Pow(1+monthRate, month) / (math.Pow(1+monthRate, month) - 1)
	fmt.Printf("等额本息(%s贷款%d年): 每个月还款%.2f\n", FormatW(principal), int(year), monthRepay)
	return
}

// FormatW
func FormatW(v float64) string {
	return fmt.Sprintf("%.2fw", v/1e4)
}

// 计算买一套房子总价为X，可用现金为Y,购买该房子付出Z%的情况下，还剩余的钱可以支撑还多少个月的房贷
func CalcMonth(total float64, cash float64, rate int) float64 {
	ratio := float64(rate*10) / 100
	loan := total * (1 - ratio)
	remainCash := cash - total*ratio
	monthLoan := EqualInterest(loan, 30)
	months := remainCash / monthLoan
	fmt.Printf("房屋总价为%s,可用现金为%s,购买该房子付出%d%%的情况下:\n还剩余的钱可以支撑还%d个月的房贷\n", FormatW(total), FormatW(cash), int(rate*10), int(months))
	return months
}
