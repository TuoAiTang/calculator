package rent

func rent(monthCost, discountAmount float64) float64 {
	return monthCost*(1+0.1*0.81)*12 - discountAmount
}
