package finance

import (
	"testing"
)

func TestParams_Calculate(t *testing.T) {
	p := &Params{
		CurrentAge:              25,
		DepositInitial:          1000000,
		CurrentMonthlyDeposit:   25000,
		YearlyDepositGrowthRate: 5,
		YearCost:                120000,
		Inflation:               3,
		FinancialIncomeRate:     6,
	}

	Output(p.Calculate())
}
