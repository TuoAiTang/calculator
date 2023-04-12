package finance

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	echarts "github.com/go-echarts/go-echarts/v2/charts"
	echarts_opts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/assert"
	"github.com/tuoaitang/calculator/db"
	"github.com/tuoaitang/calculator/model"
	"github.com/vicanso/go-charts/v2"
	"github.com/wcharczuk/go-chart/v2"
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

	err := db.Finance.Exec("delete from  yearly_stats").Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.Finance.AutoMigrate(&model.YearlyStats{})
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
				m.Inflation = p.Inflation
				m.IncomeGrowth = p.YearlyDepositGrowthRate
				models = append(models, m)
			}

			err := db.Finance.CreateInBatches(models, 1000).Error
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestParams_Calculate1(t *testing.T) {
	p := &Params{
		CurrentAge:              25,
		DepositInitial:          1000000,
		CurrentMonthlyDeposit:   25000,
		YearlyDepositGrowthRate: 5,
		YearCost:                120000,
		Inflation:               3,
		FinancialIncomeRate:     6,
	}

	stats, _ := p.Calculate()
	var xValues []string
	yValues := make([][]float64, 2)

	costValues := make([]float64, 0)
	incomeValues := make([]float64, 0)
	for _, s := range stats {
		xValues = append(xValues, strconv.FormatInt(int64(s.Age), 10))
		costValues = append(costValues, s.Cost)
		incomeValues = append(incomeValues, s.FinancialIncome)
	}

	yValues[0] = costValues
	yValues[1] = incomeValues

	painter, err := charts.LineRender(
		yValues,
		charts.XAxisOptionFunc(charts.XAxisOption{
			Data: xValues,
			//SplitNumber: 1,
		}),
		charts.TitleTextOptionFunc("income-cost"),
		charts.LegendLabelsOptionFunc([]string{"cost", "income"}, charts.PositionCenter),
		charts.ThemeOptionFunc("grafana"),
		charts.WidthOptionFunc(1000),
		charts.HeightOptionFunc(1000),
		charts.SVGTypeOption(),

		charts.YAxisOptionFunc(charts.YAxisOption{
			//Min:         thrift.Float64Ptr(12.0),
			//Max:         thrift.Float64Ptr(1.0),
			Font:        nil,
			Data:        nil,
			Theme:       nil,
			FontSize:    0,
			Position:    "left",
			FontColor:   charts.Color{},
			Formatter:   "{value}w",
			Color:       charts.Color{},
			Show:        thrift.BoolPtr(true),
			DivideCount: 10000,
			Unit:        10000,
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create(fmt.Sprintf("%d.svg", time.Now().Unix()))
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := painter.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	_, err = file.Write(bytes)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

}

func TestParams_Calculate2(t *testing.T) {
	p := &Params{
		CurrentAge:              25,
		DepositInitial:          1000000,
		CurrentMonthlyDeposit:   25000,
		YearlyDepositGrowthRate: 5,
		YearCost:                120000,
		Inflation:               3,
		FinancialIncomeRate:     6,
	}

	stats, _ := p.Calculate()
	var xValues []string

	costValues := make([]float64, 0)
	incomeValues := make([]float64, 0)
	for _, s := range stats {
		xValues = append(xValues, strconv.FormatInt(int64(s.Age), 10))
		costValues = append(costValues, s.Cost)
		incomeValues = append(incomeValues, s.FinancialIncome)
	}

	options := &charts.EChartsOption{
		Type:       "Line",
		Theme:      "",
		FontFamily: "",
		Padding:    charts.EChartsPadding{},
		Box:        chart.Box{},
		Width:      0,
		Height:     0,
		Title: struct {
			Text         string                  `json:"text"`
			Subtext      string                  `json:"subtext"`
			Left         charts.EChartsPosition  `json:"left"`
			Top          charts.EChartsPosition  `json:"top"`
			TextStyle    charts.EChartsTextStyle `json:"textStyle"`
			SubtextStyle charts.EChartsTextStyle `json:"subtextStyle"`
		}{
			Text:    "income-cost",
			Subtext: "sub-text",
		},
		XAxis: charts.EChartsXAxis{
			Data: []charts.EChartsXAxisData{
				charts.EChartsXAxisData{
					Data: xValues,
				},
			},
		},
		YAxis:  charts.EChartsYAxis{},
		Legend: charts.EChartsLegend{},
		Radar: struct {
			Indicator []charts.RadarIndicator `json:"indicator"`
		}{},
		Series: charts.EChartsSeriesList([]charts.EChartsSeries{
			charts.EChartsSeries{
				Data: []charts.EChartsSeriesData{
					charts.EChartsSeriesData{
						Value: charts.NewEChartsSeriesDataValue(costValues...),
					},
				},
			},
		}),
		Children: nil,
	}

	optionsBytes, _ := json.Marshal(options)
	bytes, err := charts.RenderEChartsToSVG(string(optionsBytes))
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create(fmt.Sprintf("%d.svg", time.Now().Unix()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = file.Write(bytes)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
}

func TestParams_Calculate3(t *testing.T) {
	p := &Params{
		CurrentAge:              25,
		DepositInitial:          1000000,
		CurrentMonthlyDeposit:   25000,
		YearlyDepositGrowthRate: 5,
		YearCost:                120000,
		Inflation:               3,
		FinancialIncomeRate:     6,
	}

	stats, _ := p.Calculate()
	var xValues []string

	costValues := make([]echarts_opts.BarData, 0)
	incomeValues := make([]echarts_opts.BarData, 0)
	for _, s := range stats {
		xValues = append(xValues, strconv.FormatInt(int64(s.Age), 10))
		costValues = append(costValues, echarts_opts.BarData{
			Value: s.Cost,
			Label: &echarts_opts.Label{
				Show:      true,
				Formatter: fmt.Sprintf("{c}w"),
			},
		})
		incomeValues = append(incomeValues, echarts_opts.BarData{
			Value: s.FinancialIncome,
			Label: &echarts_opts.Label{
				Show:      true,
				Formatter: fmt.Sprintf("{c}w"),
			},
		})
	}

	bar := echarts.NewBar()
	bar.SetGlobalOptions(
		echarts.WithTitleOpts(echarts_opts.Title{
			Title:    "年龄-收入、支出",
			Subtitle: "It's extremely easy to use, right?",
		}),
		echarts.WithLegendOpts(echarts_opts.Legend{
			Show: true,
			Data: []string{"收入", "支出"},
		}),
	)

	bar.SetSeriesOptions(
		echarts.WithLabelOpts(echarts_opts.Label{
			Show:      true,
			Color:     "",
			Position:  "",
			Formatter: "",
		}),
		echarts.WithBarChartOpts(echarts_opts.BarChart{
			Type:           "",
			Stack:          "",
			BarGap:         "",
			BarCategoryGap: "",
			XAxisIndex:     0,
			YAxisIndex:     0,
			ShowBackground: false,
			RoundCap:       false,
			CoordSystem:    "",
		}),
	)

	// Put data into instance
	bar.SetXAxis(xValues).
		AddSeries("income", incomeValues).
		AddSeries("cost", costValues)
	// Where the magic happens
	f, _ := os.Create(fmt.Sprintf("bar.html"))
	bar.Render(f)
}

func TestToLine(t *testing.T) {
	p := &Params{
		CurrentAge:              25,
		DepositInitial:          1000000,
		CurrentMonthlyDeposit:   25000,
		YearlyDepositGrowthRate: 5,
		YearCost:                120000,
		Inflation:               3,
		FinancialIncomeRate:     6,
	}

	stats, _ := p.Calculate()
	err := ToLine(stats)
	assert.Nil(t, err)
}
