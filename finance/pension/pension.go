package pension

import (
	"fmt"
	"math"
	"strings"

	"github.com/TuoAiTang/gotable/color"
	"github.com/spf13/cobra"
	"github.com/tuoaitang/calculator/util"
)

const (
	ChengduAvgSalaryYear = 96413
	AvgLifeYear          = 78.2
	BeijingUpLimit       = 31884
)

func CalcPension(cmd *cobra.Command, args []string) {
	if len(args) != 3 {
		fmt.Println("参数错误")
		return
	}

	p := Param{
		RetireCitySalaryAvg:   ChengduAvgSalaryYear / 12,
		PayYear:               parseFloat(args[0]),
		SalaryAvgBeforeRetire: parseFloat(args[1]),
		ExpectAge:             parseFloat(args[2]),
		LifeYearAvg:           AvgLifeYear,
		PayUpLimit:            BeijingUpLimit,
	}

	r := Calc(p)
	fmt.Printf("参数:\n%s\n%s\n", p, strings.Repeat("=", 50))
	fmt.Printf("结果:\n%s\n%s\n", r, strings.Repeat("=", 50))
}

func parseFloat(s string) float64 {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	if err != nil {
		panic(err)
	}
	return f
}

type Param struct {
	RetireCitySalaryAvg   float64 // 办理退休时所在城市月平均工资
	PayYear               float64 // 缴费年限
	SalaryAvgBeforeRetire float64 // 缴费期间月平均工资
	ExpectAge             float64 // 期望退休年龄
	LifeYearAvg           float64 // 平均寿命
	PayUpLimit            float64 // 缴纳所在地缴费基数上限
}

func (p Param) String() string {
	s := fmt.Sprintf("办理退休时所在城市月平均工资: %.2f\n", p.RetireCitySalaryAvg)
	s += fmt.Sprintf("期望缴费年限: %.2f\n", p.PayYear)
	s += fmt.Sprintf("缴费期间月平均工资: %.2f\n", p.SalaryAvgBeforeRetire)
	s += fmt.Sprintf("期望退休年龄: %d\n", int(p.ExpectAge))
	s += fmt.Sprintf("平均寿命: %.1f\n", p.LifeYearAvg)
	s += fmt.Sprintf("缴纳所在地缴费基数上限: %.2f\n", p.PayUpLimit)
	return s
}

type Result struct {
	// 基础养老金
	PrimaryPension float64
	// 个人账户养老金
	PersonalPension float64
	// 退休时的月平均养老金
	Pension float64
	// 退休时的月平均养老金与办理退休时所在城市月平均工资的比值
	AvgIndex float64
	// 退休时的月平均养老金与缴费期间月平均工资的比值
	BeforeRetireIndex float64
	// 领取年限
	GetYear float64
	// 总领取金额
	TotalGet float64
	// 总缴纳金额
	TotalPay float64
	// 总缴纳金额个人
	PayPerson float64
	// 总缴纳金额单位
	PayCompany float64
	// 总领取/总缴纳(个人)
	GetPayPersonIndex float64
	// 总领取/总缴纳
	GetPayIndex float64
}

func (r Result) String() string {
	s := fmt.Sprintf("领取年限: %.1f\n", r.GetYear)
	s += util.Color(color.CYAN, fmt.Sprintf("月平均养老金: %.2f(基础：%.2f + 个人账户：%.2f)\n", r.Pension, r.PrimaryPension, r.PersonalPension))
	s += fmt.Sprintf("月平均养老金与退休时所在城市月平均工资的比值: %.1f%%\n", r.AvgIndex*100)
	s += fmt.Sprintf("月平均养老金与退休前月平均工资的比值: %.1f%%\n", r.BeforeRetireIndex*100)
	s += fmt.Sprintf("总领取金额: %.2f\n", r.TotalGet)
	s += fmt.Sprintf("总缴纳金额: %.2f(个人：%.2f + 单位：%.2f)\n", r.TotalPay, r.PayPerson, r.PayCompany)
	s += util.Color(color.RED, fmt.Sprintf("总领取/总缴纳(个人): %.1f%%\n", r.GetPayPersonIndex*100))
	s += util.Color(color.RED, fmt.Sprintf("总领取/总缴纳: %.1f%%\n", r.GetPayIndex*100))
	return s
}

func Calc(p Param) Result {
	r := Result{}
	r.PrimaryPension = p.RetireCitySalaryAvg * (1 + math.Min(3, p.SalaryAvgBeforeRetire/p.RetireCitySalaryAvg)) / 2 * p.PayYear * 0.01
	r.PersonalPension = math.Min(p.SalaryAvgBeforeRetire, p.PayUpLimit) * 0.08 * p.PayYear / (p.LifeYearAvg - p.ExpectAge)
	r.Pension = r.PrimaryPension + r.PersonalPension
	r.AvgIndex = r.Pension / p.RetireCitySalaryAvg
	r.BeforeRetireIndex = r.Pension / p.SalaryAvgBeforeRetire
	r.GetYear = p.LifeYearAvg - p.ExpectAge
	r.TotalGet = r.Pension * r.GetYear * 12
	r.PayPerson = math.Min(p.SalaryAvgBeforeRetire, p.PayUpLimit) * 0.08 * 12 * p.PayYear
	r.PayCompany = math.Min(p.SalaryAvgBeforeRetire, p.PayUpLimit) * 0.20 * 12 * p.PayYear
	r.TotalPay = r.PayPerson + r.PayCompany
	r.GetPayPersonIndex = r.TotalGet / r.PayPerson
	r.GetPayIndex = r.TotalGet / r.TotalPay
	return r
}
