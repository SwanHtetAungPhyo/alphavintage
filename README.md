# Alpha Vantage Go Client

A comprehensive Go library for the Alpha Vantage API using resty.

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

    // Market Status
    status, _ := client.GetMarketStatus()
    for _, m := range status.Markets {
        fmt.Printf("%s: %s\n", m.Region, m.CurrentStatus)
    }

    // Daily Time Series
    daily, _ := client.GetTimeSeriesDaily("IBM", alphavintage.OutputSizeCompact)
    fmt.Printf("Symbol: %s\n", daily.MetaData.Symbol)

    // Intraday Time Series
    intraday, _ := client.GetTimeSeriesIntraday("IBM", alphavintage.Interval5Min, alphavintage.OutputSizeCompact)
    fmt.Printf("Last: %s\n", intraday.MetaData.LastRefreshed)

    // Balance Sheet
    balance, _ := client.GetBalanceSheet("IBM")
    fmt.Printf("Total Assets: %s\n", balance.AnnualReports[0].TotalAssets)

    // Cash Flow
    cashflow, _ := client.GetCashFlow("IBM")
    fmt.Printf("Net Income: %s\n", cashflow.AnnualReports[0].NetIncome)

    // Earnings
    earnings, _ := client.GetEarnings("IBM")
    fmt.Printf("EPS: %s\n", earnings.AnnualEarnings[0].ReportedEPS)

    // News Sentiment
    news, _ := client.GetNewsSentiment(&alphavintage.NewsSentimentOptions{
        Tickers: "AAPL",
    })
    fmt.Printf("Articles: %d\n", len(news.Feed))
}
```

## API Functions

| Function | Description |
|----------|-------------|
| `GetMarketStatus()` | Global market open/close status |
| `GetTimeSeriesDaily(symbol, outputSize)` | Daily OHLCV data |
| `GetTimeSeriesIntraday(symbol, interval, outputSize)` | Intraday OHLCV data |
| `GetBalanceSheet(symbol)` | Balance sheet fundamentals |
| `GetCashFlow(symbol)` | Cash flow statements |
| `GetEarnings(symbol)` | Earnings data |
| `GetNewsSentiment(options)` | News and sentiment analysis |

## Options

### OutputSize
- `OutputSizeCompact` - Last 100 data points
- `OutputSizeFull` - Full historical data

### Interval (Intraday)
- `Interval1Min`, `Interval5Min`, `Interval15Min`, `Interval30Min`, `Interval60Min`

## License

MIT
