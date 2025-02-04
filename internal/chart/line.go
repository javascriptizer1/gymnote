package chart

import (
	"log"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/snapshot-chromedp/render"
)

type LinearChartConfig struct {
	Title    string
	XName    string
	YName    string
	YValues  []float32
	XValues  []string
	FileName string
}

func (c *chart) GenerateLinearChart(config LinearChartConfig) error {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: config.Title}),
		charts.WithAnimation(false),
		charts.WithXAxisOpts(opts.XAxis{Name: config.XName}),
		charts.WithYAxisOpts(opts.YAxis{Name: config.YName}),
	)

	var lineData []opts.LineData
	for _, value := range config.YValues {
		lineData = append(lineData, opts.LineData{Value: value})
	}

	line.SetXAxis(config.XValues).AddSeries(config.YName, lineData).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show: opts.Bool(true),
			}),
			charts.WithAreaStyleOpts(opts.AreaStyle{
				Opacity: 0.2,
			}),
			charts.WithLineChartOpts(opts.LineChart{
				Smooth: opts.Bool(true),
			}),
		)

	if err := render.MakeSnapshot(render.NewSnapshotConfig(line.RenderContent(), config.FileName, func(config *render.SnapshotConfig) {
		config.MultiCharts = true
		config.KeepHtml = true
		config.Quality = 100
	})); err != nil {
		log.Printf("Render chart error: %v\n", err)
		return err
	}

	return nil
}
