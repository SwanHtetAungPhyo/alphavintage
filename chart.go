package alphavintage

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

// ChartOptions configures chart generation
type ChartOptions struct {
	Width      int
	Height     int
	Title      string
	ShowVolume bool
}

// DefaultChartOptions returns default chart options
func DefaultChartOptions() ChartOptions {
	return ChartOptions{
		Width:      1200,
		Height:     600,
		ShowVolume: true,
	}
}

// GenerateDailyPriceChart creates a price chart from daily time series data
func GenerateDailyPriceChart(data *TimeSeriesDailyResponse, output io.Writer, opts ChartOptions) error {
	if data == nil || len(data.TimeSeries) == 0 {
		return fmt.Errorf("no data to chart")
	}

	if opts.Width == 0 {
		opts.Width = 1200
	}
	if opts.Height == 0 {
		opts.Height = 600
	}
	if opts.Title == "" {
		opts.Title = fmt.Sprintf("%s Daily Price", data.MetaData.Symbol)
	}

	// Sort dates and extract data
	dates, closes, volumes := extractDailyData(data.TimeSeries)

	// Create price series
	priceSeries := chart.TimeSeries{
		Name: "Close Price",
		Style: chart.Style{
			StrokeColor: chart.ColorBlue,
			StrokeWidth: 2,
		},
		XValues: dates,
		YValues: closes,
	}

	graph := chart.Chart{
		Title:      opts.Title,
		TitleStyle: chart.Style{FontSize: 14},
		Width:      opts.Width,
		Height:     opts.Height,
		XAxis: chart.XAxis{
			Name:           "Date",
			TickPosition:   chart.TickPositionBetweenTicks,
			ValueFormatter: chart.TimeDateValueFormatter,
		},
		YAxis: chart.YAxis{
			Name: "Price ($)",
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("$%.2f", v.(float64))
			},
		},
		Series: []chart.Series{priceSeries},
	}

	// Add volume bars if requested
	if opts.ShowVolume && len(volumes) > 0 {
		graph.YAxisSecondary = chart.YAxis{
			Name: "Volume",
			ValueFormatter: func(v interface{}) string {
				return formatVolume(v.(float64))
			},
		}
		volumeSeries := chart.TimeSeries{
			Name:    "Volume",
			YAxis:   chart.YAxisSecondary,
			XValues: dates,
			YValues: volumes,
			Style: chart.Style{
				StrokeColor: drawing.ColorFromHex("90EE90"),
				FillColor:   drawing.ColorFromHex("90EE90").WithAlpha(100),
				StrokeWidth: 0,
			},
		}
		graph.Series = append(graph.Series, volumeSeries)
	}

	graph.Elements = []chart.Renderable{chart.Legend(&graph)}

	return graph.Render(chart.PNG, output)
}

// GenerateDailyPriceChartToFile saves chart to a PNG file
func GenerateDailyPriceChartToFile(data *TimeSeriesDailyResponse, filename string, opts ChartOptions) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return GenerateDailyPriceChart(data, f, opts)
}


// GenerateCandlestickChart creates a candlestick chart from daily data
func GenerateCandlestickChart(data *TimeSeriesDailyResponse, output io.Writer, opts ChartOptions) error {
	if data == nil || len(data.TimeSeries) == 0 {
		return fmt.Errorf("no data to chart")
	}

	if opts.Width == 0 {
		opts.Width = 1200
	}
	if opts.Height == 0 {
		opts.Height = 600
	}
	if opts.Title == "" {
		opts.Title = fmt.Sprintf("%s Candlestick Chart", data.MetaData.Symbol)
	}

	// Sort and extract OHLC data
	type ohlc struct {
		date                       time.Time
		open, high, low, close     float64
	}

	var candles []ohlc
	for dateStr, dp := range data.TimeSeries {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		o, _ := strconv.ParseFloat(dp.Open, 64)
		h, _ := strconv.ParseFloat(dp.High, 64)
		l, _ := strconv.ParseFloat(dp.Low, 64)
		c, _ := strconv.ParseFloat(dp.Close, 64)
		candles = append(candles, ohlc{t, o, h, l, c})
	}

	sort.Slice(candles, func(i, j int) bool {
		return candles[i].date.Before(candles[j].date)
	})

	// Create high/low range and close line
	var dates []time.Time
	var highs, lows, closes []float64

	for _, c := range candles {
		dates = append(dates, c.date)
		highs = append(highs, c.high)
		lows = append(lows, c.low)
		closes = append(closes, c.close)
	}

	highSeries := chart.TimeSeries{
		Name:    "High",
		XValues: dates,
		YValues: highs,
		Style: chart.Style{
			StrokeColor: drawing.ColorFromHex("28a745"),
			StrokeWidth: 1,
			DotWidth:    2,
		},
	}

	lowSeries := chart.TimeSeries{
		Name:    "Low",
		XValues: dates,
		YValues: lows,
		Style: chart.Style{
			StrokeColor: drawing.ColorFromHex("dc3545"),
			StrokeWidth: 1,
			DotWidth:    2,
		},
	}

	closeSeries := chart.TimeSeries{
		Name:    "Close",
		XValues: dates,
		YValues: closes,
		Style: chart.Style{
			StrokeColor: chart.ColorBlue,
			StrokeWidth: 2,
		},
	}

	graph := chart.Chart{
		Title:      opts.Title,
		TitleStyle: chart.Style{FontSize: 14},
		Width:      opts.Width,
		Height:     opts.Height,
		XAxis: chart.XAxis{
			Name:           "Date",
			ValueFormatter: chart.TimeDateValueFormatter,
		},
		YAxis: chart.YAxis{
			Name: "Price ($)",
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("$%.2f", v.(float64))
			},
		},
		Series: []chart.Series{highSeries, lowSeries, closeSeries},
	}

	graph.Elements = []chart.Renderable{chart.Legend(&graph)}

	return graph.Render(chart.PNG, output)
}

// GenerateCandlestickChartToFile saves candlestick chart to PNG file
func GenerateCandlestickChartToFile(data *TimeSeriesDailyResponse, filename string, opts ChartOptions) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return GenerateCandlestickChart(data, f, opts)
}


// GenerateEarningsChart creates a bar chart of earnings over time
func GenerateEarningsChart(data *EarningsResponse, output io.Writer, opts ChartOptions) error {
	if data == nil || len(data.AnnualEarnings) == 0 {
		return fmt.Errorf("no earnings data to chart")
	}

	if opts.Width == 0 {
		opts.Width = 800
	}
	if opts.Height == 0 {
		opts.Height = 400
	}
	if opts.Title == "" {
		opts.Title = fmt.Sprintf("%s Annual EPS", data.Symbol)
	}

	// Sort by date and limit to recent years
	type earning struct {
		date time.Time
		eps  float64
	}

	var earnings []earning
	for _, e := range data.AnnualEarnings {
		t, err := time.Parse("2006-01-02", e.FiscalDateEnding)
		if err != nil {
			continue
		}
		eps, _ := strconv.ParseFloat(e.ReportedEPS, 64)
		earnings = append(earnings, earning{t, eps})
	}

	sort.Slice(earnings, func(i, j int) bool {
		return earnings[i].date.Before(earnings[j].date)
	})

	// Limit to last 10 years
	if len(earnings) > 10 {
		earnings = earnings[len(earnings)-10:]
	}

	var bars []chart.Value
	for _, e := range earnings {
		bars = append(bars, chart.Value{
			Label: e.date.Format("2006"),
			Value: e.eps,
		})
	}

	graph := chart.BarChart{
		Title:      opts.Title,
		TitleStyle: chart.Style{FontSize: 14},
		Width:      opts.Width,
		Height:     opts.Height,
		BarWidth:   40,
		XAxis: chart.Style{
			FontSize: 10,
		},
		YAxis: chart.YAxis{
			Name: "EPS ($)",
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("$%.2f", v.(float64))
			},
		},
		Bars: bars,
	}

	return graph.Render(chart.PNG, output)
}

// GenerateEarningsChartToFile saves earnings chart to PNG file
func GenerateEarningsChartToFile(data *EarningsResponse, filename string, opts ChartOptions) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return GenerateEarningsChart(data, f, opts)
}

// GenerateComparisonChart creates a multi-line chart comparing multiple symbols
func GenerateComparisonChart(datasets map[string]*TimeSeriesDailyResponse, output io.Writer, opts ChartOptions) error {
	if len(datasets) == 0 {
		return fmt.Errorf("no data to chart")
	}

	if opts.Width == 0 {
		opts.Width = 1200
	}
	if opts.Height == 0 {
		opts.Height = 600
	}
	if opts.Title == "" {
		opts.Title = "Price Comparison"
	}

	colors := []drawing.Color{
		chart.ColorBlue,
		chart.ColorRed,
		chart.ColorGreen,
		chart.ColorOrange,
		chart.ColorCyan,
	}

	var series []chart.Series
	colorIdx := 0

	for symbol, data := range datasets {
		if data == nil || len(data.TimeSeries) == 0 {
			continue
		}

		dates, closes, _ := extractDailyData(data.TimeSeries)

		// Normalize to percentage change from first value
		if len(closes) > 0 {
			base := closes[0]
			normalized := make([]float64, len(closes))
			for i, v := range closes {
				normalized[i] = ((v - base) / base) * 100
			}

			series = append(series, chart.TimeSeries{
				Name:    symbol,
				XValues: dates,
				YValues: normalized,
				Style: chart.Style{
					StrokeColor: colors[colorIdx%len(colors)],
					StrokeWidth: 2,
				},
			})
			colorIdx++
		}
	}

	if len(series) == 0 {
		return fmt.Errorf("no valid data to chart")
	}

	graph := chart.Chart{
		Title:      opts.Title,
		TitleStyle: chart.Style{FontSize: 14},
		Width:      opts.Width,
		Height:     opts.Height,
		XAxis: chart.XAxis{
			Name:           "Date",
			ValueFormatter: chart.TimeDateValueFormatter,
		},
		YAxis: chart.YAxis{
			Name: "Change (%)",
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%.1f%%", v.(float64))
			},
		},
		Series: series,
	}

	graph.Elements = []chart.Renderable{chart.Legend(&graph)}

	return graph.Render(chart.PNG, output)
}

// GenerateComparisonChartToFile saves comparison chart to PNG file
func GenerateComparisonChartToFile(datasets map[string]*TimeSeriesDailyResponse, filename string, opts ChartOptions) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return GenerateComparisonChart(datasets, f, opts)
}


// Helper functions

func extractDailyData(timeSeries map[string]DailyDataPoint) ([]time.Time, []float64, []float64) {
	type dataPoint struct {
		date   time.Time
		close  float64
		volume float64
	}

	var points []dataPoint
	for dateStr, dp := range timeSeries {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		close, _ := strconv.ParseFloat(dp.Close, 64)
		vol, _ := strconv.ParseFloat(dp.Volume, 64)
		points = append(points, dataPoint{t, close, vol})
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].date.Before(points[j].date)
	})

	var dates []time.Time
	var closes, volumes []float64
	for _, p := range points {
		dates = append(dates, p.date)
		closes = append(closes, p.close)
		volumes = append(volumes, p.volume)
	}

	return dates, closes, volumes
}

func formatVolume(v float64) string {
	if v >= 1e9 {
		return fmt.Sprintf("%.1fB", v/1e9)
	}
	if v >= 1e6 {
		return fmt.Sprintf("%.1fM", v/1e6)
	}
	if v >= 1e3 {
		return fmt.Sprintf("%.1fK", v/1e3)
	}
	return fmt.Sprintf("%.0f", v)
}

// GenerateCashFlowChart creates a chart showing cash flow trends
func GenerateCashFlowChart(data *CashFlowResponse, output io.Writer, opts ChartOptions) error {
	if data == nil || len(data.AnnualReports) == 0 {
		return fmt.Errorf("no cash flow data to chart")
	}

	if opts.Width == 0 {
		opts.Width = 1000
	}
	if opts.Height == 0 {
		opts.Height = 500
	}
	if opts.Title == "" {
		opts.Title = fmt.Sprintf("%s Cash Flow", data.Symbol)
	}

	type cfData struct {
		date       time.Time
		operating  float64
		investing  float64
		financing  float64
	}

	var cfPoints []cfData
	for _, r := range data.AnnualReports {
		t, err := time.Parse("2006-01-02", r.FiscalDateEnding)
		if err != nil {
			continue
		}
		op, _ := strconv.ParseFloat(r.OperatingCashflow, 64)
		inv, _ := strconv.ParseFloat(r.CashflowFromInvestment, 64)
		fin, _ := strconv.ParseFloat(r.CashflowFromFinancing, 64)
		cfPoints = append(cfPoints, cfData{t, op / 1e9, inv / 1e9, fin / 1e9})
	}

	sort.Slice(cfPoints, func(i, j int) bool {
		return cfPoints[i].date.Before(cfPoints[j].date)
	})

	if len(cfPoints) > 10 {
		cfPoints = cfPoints[len(cfPoints)-10:]
	}

	var dates []time.Time
	var operating, investing, financing []float64
	for _, p := range cfPoints {
		dates = append(dates, p.date)
		operating = append(operating, p.operating)
		investing = append(investing, p.investing)
		financing = append(financing, p.financing)
	}

	graph := chart.Chart{
		Title:      opts.Title,
		TitleStyle: chart.Style{FontSize: 14},
		Width:      opts.Width,
		Height:     opts.Height,
		XAxis: chart.XAxis{
			Name:           "Year",
			ValueFormatter: chart.TimeDateValueFormatter,
		},
		YAxis: chart.YAxis{
			Name: "Cash Flow ($B)",
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("$%.1fB", v.(float64))
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "Operating",
				XValues: dates,
				YValues: operating,
				Style: chart.Style{
					StrokeColor: chart.ColorGreen,
					StrokeWidth: 2,
				},
			},
			chart.TimeSeries{
				Name:    "Investing",
				XValues: dates,
				YValues: investing,
				Style: chart.Style{
					StrokeColor: chart.ColorBlue,
					StrokeWidth: 2,
				},
			},
			chart.TimeSeries{
				Name:    "Financing",
				XValues: dates,
				YValues: financing,
				Style: chart.Style{
					StrokeColor: chart.ColorRed,
					StrokeWidth: 2,
				},
			},
		},
	}

	graph.Elements = []chart.Renderable{chart.Legend(&graph)}

	return graph.Render(chart.PNG, output)
}

// GenerateCashFlowChartToFile saves cash flow chart to PNG file
func GenerateCashFlowChartToFile(data *CashFlowResponse, filename string, opts ChartOptions) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return GenerateCashFlowChart(data, f, opts)
}
