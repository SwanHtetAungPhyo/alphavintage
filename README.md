# Alpha Vantage Go Client

A comprehensive Go library for the Alpha Vantage API with charting support.

## Installation

```bash
go get github.com/SwanHtetAungPhyo/alphavintage
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/SwanHtetAungPhyo/alphavintage"
)

func main() {
    client := alphavintage.NewClient("YOUR_API_KEY")

    // Get daily time series
    daily, _ := client.GetTimeSeriesDaily("IBM", alphavintage.OutputSizeCompact)
    
    // Generate price chart PNG
    opts := alphavintage.DefaultChartOptions()
    opts.Title = "IBM Stock Price"
    alphavintage.GenerateDailyPriceChartToFile(daily, "price.png", opts)
    
    // Generate candlestick chart
    alphavintage.GenerateCandlestickChartToFile(daily, "candlestick.png", opts)
    
    // Get and chart earnings
    earnings, _ := client.GetEarnings("IBM")
    alphavintage.GenerateEarningsChartToFile(earnings, "earnings.png", opts)
    
    // Get and chart cash flow
    cashflow, _ := client.GetCashFlow("IBM")
    alphavintage.GenerateCashFlowChartToFile(cashflow, "cashflow.png", opts)
}
```

## API Functions

| Function | Description |
|----------|-------------|
| `GetMarketStatus()` | Global market open/close status |
| `GetTimeSeriesDaily(symbol, outputSize)` | Daily OHLCV data |
| `GetTimeSeriesIntraday(symbol, interval, outputSize)` | Intraday OHLCV (premium) |
| `GetBalanceSheet(symbol)` | Balance sheet fundamentals |
| `GetCashFlow(symbol)` | Cash flow statements |
| `GetEarnings(symbol)` | Earnings data |
| `GetNewsSentiment(options)` | News and sentiment analysis |

## Chart Functions

| Function | Description |
|----------|-------------|
| `GenerateDailyPriceChart(data, writer, opts)` | Line chart with volume |
| `GenerateDailyPriceChartToFile(data, filename, opts)` | Save price chart to PNG |
| `GenerateCandlestickChart(data, writer, opts)` | High/Low/Close chart |
| `GenerateCandlestickChartToFile(data, filename, opts)` | Save candlestick to PNG |
| `GenerateEarningsChart(data, writer, opts)` | Annual EPS bar chart |
| `GenerateEarningsChartToFile(data, filename, opts)` | Save earnings to PNG |
| `GenerateCashFlowChart(data, writer, opts)` | Cash flow trends |
| `GenerateCashFlowChartToFile(data, filename, opts)` | Save cash flow to PNG |
| `GenerateComparisonChart(datasets, writer, opts)` | Multi-symbol comparison |
| `GenerateComparisonChartToFile(datasets, filename, opts)` | Save comparison to PNG |

## Chart Options

```go
opts := alphavintage.ChartOptions{
    Width:      1200,    // Chart width in pixels
    Height:     600,     // Chart height in pixels
    Title:      "Title", // Chart title
    ShowVolume: true,    // Show volume bars (price chart only)
}
```

## Using Charts in PDFs

The generated PNG files can be embedded in PDFs using libraries like:
- `github.com/jung-kurt/gofpdf`
- `github.com/signintech/gopdf`

```go
import "github.com/jung-kurt/gofpdf"

pdf := gofpdf.New("P", "mm", "A4", "")
pdf.AddPage()
pdf.Image("ibm_price.png", 10, 10, 190, 0, false, "", 0, "")
pdf.OutputFileAndClose("report.pdf")
```

## Options

### OutputSize
- `OutputSizeCompact` - Last 100 data points
- `OutputSizeFull` - Full historical data

### Interval (Intraday - Premium)
- `Interval1Min`, `Interval5Min`, `Interval15Min`, `Interval30Min`, `Interval60Min`

## License

MIT
