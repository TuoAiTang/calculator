package finance

import (
	"fmt"
)

type Params struct {
	CurrentAge              int     // 当前年龄
	DepositInitial          float64 // 初始存款
	CurrentMonthlyDeposit   float64 // 当前每月存款
	YearlyDepositGrowthRate float64 // 年存款增长率 * 100
	YearCost                float64 // 每年花费
	Inflation               float64 // 通货膨胀率 * 100
	FinancialIncomeRate     float64 // 理财收益率 * 100
}

// Print 输出
func (p *Params) Print() {
	fmt.Printf("当前年龄: %d\n", p.CurrentAge)
	fmt.Printf("初始存款: %.2f\n", p.DepositInitial)
	fmt.Printf("当前每月存款: %.2f\n", p.CurrentMonthlyDeposit)
	fmt.Printf("年存款增长率: %.2f%%\n", p.YearlyDepositGrowthRate)
	fmt.Printf("每年花费: %.2f\n", p.YearCost)
	fmt.Printf("通货膨胀率: %.2f%%\n", p.Inflation)
	fmt.Printf("理财收益率: %.2f%%\n", p.FinancialIncomeRate)
}

// YearlyStats 每年的统计数据
type YearlyStats struct {
	Age             int     // 年龄
	YearEndDeposit  float64 // 年末存款
	Cost            float64 // 支出
	FinancialIncome float64 // 理财收入
	YearlyDeposit   float64 // 年存款
}

// Calculate 计算
func (p *Params) Calculate() []YearlyStats {
	p.Print()
	yearlyStats := make([]YearlyStats, 0)
	yearlyStats = append(yearlyStats, YearlyStats{
		Age:             p.CurrentAge,
		YearEndDeposit:  p.DepositInitial,
		Cost:            p.YearCost,
		FinancialIncome: 0,
		YearlyDeposit:   p.CurrentMonthlyDeposit * 12 / (1 + p.YearlyDepositGrowthRate/100),
	})

	lastStats := yearlyStats[0]
	for lastStats.Age < 65 {
		currentStats := YearlyStats{
			Age:             lastStats.Age + 1,
			YearEndDeposit:  lastStats.YearEndDeposit,
			Cost:            lastStats.Cost * (1 + p.Inflation/100),
			FinancialIncome: lastStats.YearEndDeposit * p.FinancialIncomeRate / 100,
			YearlyDeposit:   lastStats.YearlyDeposit * (1 + p.YearlyDepositGrowthRate/100),
		}

		currentStats.YearEndDeposit = lastStats.YearEndDeposit + currentStats.FinancialIncome + currentStats.YearlyDeposit
		yearlyStats = append(yearlyStats, currentStats)
		lastStats = currentStats
	}
	return yearlyStats
}

// Output
func Output(yearlyStats []YearlyStats) {
	for _, stats := range yearlyStats {
		fmt.Printf("年龄: %d, 年末存款: %s, 支出: %s, 理财收入: %s\n", stats.Age, FormatW(stats.YearEndDeposit), FormatW(stats.Cost), FormatW(stats.FinancialIncome))
	}
}

// FormatW
func FormatW(v float64) string {
	return fmt.Sprintf("%.2fw", v/1e4)
}
