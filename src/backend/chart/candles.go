package chart

import (
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type KlineData struct {
	Date string
	Data [4]float32
}

var ksd = []KlineData{
	{Date: "2010/1/24", Data: [4]float32{2320.26, 2320.26, 2287.3, 2362.94}},
	{Date: "2015/1/25", Data: [4]float32{2300, 2291.3, 2288.26, 2308.38}},
	{Date: "2017/1/28", Data: [4]float32{2295.35, 2346.5, 2295.35, 2346.92}},
	{Date: "2018/1/29", Data: [4]float32{2347.22, 2358.98, 2337.35, 2363.8}},
	{Date: "2018/1/30", Data: [4]float32{2360.75, 2382.48, 2347.89, 2383.76}},
}

func KlineDataZoomInside(kd []KlineData) *charts.Kline {
	kline := charts.NewKLine()

	//now := time.Now()

	// Extract the parts we want
	//year, month, day := now.Date()
	//hour, min, _ := now.Clock()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].Date)
		y = append(y, opts.KlineData{Value: kd[i].Data})
	}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "DataZoom(inside)",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
			//Type:  "time",
			//Min:   time.Date(year, month, day, hour, min, 0, 0, time.Local),
			//Max:   now.Add(time.Duration(len(kd)) * time.Minute), // for example, 1-minute intervals
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("kline", y).
		SetSeriesOptions(
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color:        "#00da3c",
				Color0:       "#ec0000",
				BorderColor:  "#00da3c",
				BorderColor0: "#ec0000",
			}),
		)
	return kline
}

type KlineExamples struct{}

func (KlineExamples) Examples(kd []KlineData) {
	page := components.NewPage()
	page.SetLayout(components.PageFullLayout)
	page.AddCharts(
		KlineDataZoomInside(kd),
	)

	f, err := os.Create("./kline.html")
	if err != nil {
		panic(err)

	}
	page.Render(io.MultiWriter(f))
}
