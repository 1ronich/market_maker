package chart

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateLineItemsTwoAxis(points int, values []float64, xFunc func(int) interface{}) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < points; i++ {
		fmt.Println(values[i])
		items = append(items, opts.LineData{Value: []interface{}{xFunc(i), values[i]}})
	}
	return items
}

func lineTime(values []float64) *charts.Line {
	now := time.Now()

	// Extract the parts we want
	year, month, day := now.Date()
	hour, min, _ := now.Clock()

	// Create a new time.Time with those values

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "temporal X axis",
			Subtitle: "time.Date as X axis values",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Min: 0.00008,
			Max: 0.0003,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "time",
			Min:  time.Date(year, month, day, hour, min, 0, 0, time.Local),
			Max:  now.Add(time.Duration(len(values)) * time.Minute), // for example, 1-minute intervals
		}),
		charts.WithTooltipOpts(opts.Tooltip{ // Potential to string format tooltip here
			Show:    opts.Bool(true),
			Trigger: "axis",
		}),
	)

	fmt.Println("len:", len(values))

	line.AddSeries("Category A", generateLineItemsTwoAxis(len(values), values, func(i int) interface{} {
		return now.Add(time.Duration(i) * time.Minute)
	}))
	return line
}

type LineExamples struct{}

func (LineExamples) Examples(values []float64) {
	page := components.NewPage()
	page.AddCharts(

		lineTime(values),
	)
	f, err := os.Create("./line.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
