package finance

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/TuoAiTang/gotable/color"
	"github.com/TuoAiTang/gotable/table"
	echarts "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/tuoaitang/calculator/model"
)

type FreedomCalculator struct {
	Params *Params
}

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
	Age             int     `gorm:"column:age"`              // 年龄
	YearEndDeposit  float64 `gorm:"column:year_end_deposit"` // 年末存款
	Cost            float64 `gorm:"column:cost"`             // 支出
	FinancialIncome float64 `gorm:"column:financial_income"` // 理财收入
	YearlyDeposit   float64 `gorm:"column:yearly_deposit"`   // 年存款
	CanCover        bool    `gorm:"column:can_cover"`        // 是否能够覆盖支出
}

// ToModel
func (y YearlyStats) ToModel() model.YearlyStats {
	return model.YearlyStats{
		Age:             y.Age,
		YearEndDeposit:  FormatW(y.YearEndDeposit),
		Cost:            FormatW(y.Cost),
		FinancialIncome: FormatW(y.FinancialIncome),
		YearlyDeposit:   FormatW(y.YearlyDeposit),
		CanCover:        y.CanCover,
	}
}

// IsFreedom
func (y YearlyStats) IsFreedom() bool {
	return y.FinancialIncome > y.Cost
}

func ToModels(stats []YearlyStats) []model.YearlyStats {
	models := make([]model.YearlyStats, 0)
	for _, s := range stats {
		models = append(models, s.ToModel())
	}
	return models
}

// Calculate 计算
func (p *Params) Calculate() ([]YearlyStats, int) {
	p.Print()
	yearlyStats := make([]YearlyStats, 0)
	yearlyStats = append(yearlyStats, YearlyStats{
		Age:             p.CurrentAge,
		YearEndDeposit:  p.DepositInitial,
		Cost:            p.YearCost,
		FinancialIncome: 0,
		YearlyDeposit:   p.CurrentMonthlyDeposit * 12 / (1 + p.YearlyDepositGrowthRate/100),
		CanCover:        false,
	})

	lastStats := yearlyStats[0]
	index := -1
	for lastStats.Age < 60 {
		currentStats := YearlyStats{
			Age:             lastStats.Age + 1,
			YearEndDeposit:  lastStats.YearEndDeposit,
			Cost:            lastStats.Cost * (1 + p.Inflation/100),
			FinancialIncome: lastStats.YearEndDeposit * p.FinancialIncomeRate / 100,
			YearlyDeposit:   lastStats.YearlyDeposit * (1 + p.YearlyDepositGrowthRate/100),
		}

		currentStats.YearEndDeposit = lastStats.YearEndDeposit + currentStats.FinancialIncome + currentStats.YearlyDeposit
		currentStats.CanCover = currentStats.FinancialIncome > currentStats.Cost
		yearlyStats = append(yearlyStats, currentStats)
		if currentStats.CanCover && index == -1 {
			index = len(yearlyStats) - 1
		}

		lastStats = currentStats
	}
	return yearlyStats, index
}

// Output
func Output(yearlyStats []YearlyStats) {
	type data struct {
		Age             string
		YearEndDeposit  string
		Cost            string
		FinancialIncome string
	}

	var firstColoredFields int
	free := func(d data) bool {
		income, _ := strconv.ParseFloat(strings.Trim(d.FinancialIncome, "w"), 64)
		cost, _ := strconv.ParseFloat(strings.Trim(d.Cost, "w"), 64)
		return income > cost
	}

	dataSlice := make([]interface{}, 0)
	//dataSlice = append(dataSlice, data{
	//	Age:             "年龄",
	//	YearEndDeposit:  "年末存款",
	//	Cost:            "支出",
	//	FinancialIncome: "理财收入",
	//})
	for _, v := range yearlyStats {
		dataSlice = append(dataSlice, data{
			Age:             strconv.FormatInt(int64(v.Age), 10),
			YearEndDeposit:  FormatW(v.YearEndDeposit),
			Cost:            FormatW(v.Cost),
			FinancialIncome: FormatW(v.FinancialIncome),
		})
	}

	//t, err := table.CreateTable([]string{"年龄", "年末存款", "支出", "理财收入"}, table.WithRecordColorController(func(record interface{}) color.Color {
	t, err := table.Create(dataSlice, table.WithRecordColorController(func(record interface{}) color.Color {
		r := record.(data)
		if free(r) {
			firstColoredFields++
			if firstColoredFields <= 4 {
				return color.RED
			}
			return color.CYAN
		}
		return ""
	}))

	if err != nil {
		panic(err)
	}

	t.PrintTable()
}

// FormatW
func FormatW(v float64) string {
	return fmt.Sprintf("%.2fw", v/1e4)
}

func ToLine(data []YearlyStats) error {
	line := echarts.NewLine()

	line.SetGlobalOptions(
		echarts.WithTitleOpts(opts.Title{
			Title:    "年龄-收入、支出",
			Subtitle: "It's extremely easy to use, right?",
		}),
	)
	line.SetSeriesOptions(
		echarts.WithLabelOpts(opts.Label{
			Show:      true,
			Color:     "",
			Position:  "",
			Formatter: "{a}",
		}),
	)
	xValues := make([]string, 0)
	incomeValues := make([]opts.LineData, 0)
	costValues := make([]opts.LineData, 0)
	for _, v := range data {
		xValues = append(xValues, strconv.FormatInt(int64(v.Age), 10))
		incomeValues = append(incomeValues, opts.LineData{Value: v.FinancialIncome})
		costValues = append(costValues, opts.LineData{Value: v.Cost})
	}
	line.SetXAxis(xValues).AddSeries("理财收入", incomeValues).AddSeries("支出", costValues)

	// Where the magic happens
	f, _ := os.Create(fmt.Sprintf("line.html"))
	line.Render(f)
	return nil
}
