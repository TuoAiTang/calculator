package db

import (
	"testing"
)

func TestOpen(t *testing.T) {
	Init("finance.db")
	//Finance.AutoMigrate(&model.YearlyStats{})
}
