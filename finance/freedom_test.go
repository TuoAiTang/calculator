package finance

import (
	"fmt"
	"testing"

	"github.com/tuoaitang/calculator/db"
	"github.com/tuoaitang/calculator/model"
)

func TestParams_Calculate(t *testing.T) {
	p := &Params{
		CurrentAge:              25,
		DepositInitial:          1000000,
		CurrentMonthlyDeposit:   25000,
		YearlyDepositGrowthRate: 5,
		YearCost:                120000,
		Inflation:               20,
		FinancialIncomeRate:     3,
	}

	err := db.Finance.AutoMigrate(&model.YearlyStats{})
	if err != nil {
		t.Fatal(err)
	}
	err = db.Finance.Exec("DELETE FROM yearly_stats").Error
	if err != nil {
		t.Fatal(err)
	}
	for inflation := 0; inflation <= 30; inflation++ {
		p.Inflation = float64(inflation)
		for incomeGrowth := 0; incomeGrowth <= 30; incomeGrowth++ {
			p.YearlyDepositGrowthRate = float64(incomeGrowth)
			var models []model.YearlyStats
			stats, _ := p.Calculate()
			for _, s := range stats {
				m := s.ToModel()
				m.Inflation = fmt.Sprintf("%.2f%%", p.Inflation)
				m.IncomeGrowth = fmt.Sprintf("%.2f%%", p.YearlyDepositGrowthRate)
				models = append(models, m)
			}

			err := db.Finance.CreateInBatches(models, 1000).Error
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
