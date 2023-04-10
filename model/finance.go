package model

// YearlyStats 每年的统计数据
type YearlyStats struct {
	Inflation       string // 通货膨胀率
	IncomeGrowth    string // 收入增长率
	Age             int    // 年龄
	YearEndDeposit  string // 年末存款
	Cost            string // 支出
	FinancialIncome string // 理财收入
	YearlyDeposit   string // 年存款
	CanCover        bool   // 是否能够覆盖支出
}
