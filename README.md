# Alpha Vantage Go Client

A comprehensive Go library for Alpha Vantage API with charting, PDF reports, and AI-powered analysis.

## Installation

```bash
go get github.com/SwanHtetAungPhyo/alphavintage
```

## Quick Start

```go
client := alphavintage.NewClient("YOUR_ALPHA_VANTAGE_KEY")
daily, _ := client.GetTimeSeriesDaily("IBM", alphavintage.OutputSizeCompact)

// Generate chart
alphavintage.GenerateDailyPriceChartToFile(daily, "price.png", alphavintage.ChartOptions{})

// Generate PDF report
report := alphavintage.NewReportBuilder(alphavintage.DefaultReportOptions())
report.AddPage().AddTitle("IBM Analysis")
report.AddDailyPriceChart(daily, alphavintage.ChartOptions{})
report.Save("report.pdf")
```

## API Functions

| Function | Description |
|----------|-------------|
| `GetMarketStatus()` | Global market status |
| `GetTimeSeriesDaily(symbol, outputSize)` | Daily OHLCV |
| `GetTimeSeriesIntraday(symbol, interval, outputSize)` | Intraday (premium) |
| `GetBalanceSheet(symbol)` | Balance sheet |
| `GetCashFlow(symbol)` | Cash flow |
| `GetEarnings(symbol)` | Earnings |
| `GetNewsSentiment(options)` | News sentiment |

## AI-Powered Analysis

Generate AI summaries using OpenRouter (supports multiple models):

```go
// Configure AI client
aiConfig := alphavintage.AIConfig{
    APIKey:    "YOUR_OPENROUTER_API_KEY",
    Model:     "nvidia/nemotron-3-nano-30b-a3b:free", // Free model
    // Model:  "openai/gpt-4o-mini",                  // Paid model
    // Model:  "anthropic/claude-3-haiku",            // Another option
    Reasoning: false,
}
aiClient := alphavintage.NewAIClient(aiConfig)

// Prepare data
stockData := alphavintage.StockAnalysisData{
    Symbol:       "IBM",
    Daily:        daily,
    Earnings:     earnings,
    CashFlow:     cashflow,
    BalanceSheet: balance,
}

// Generate full analysis
summary, _ := aiClient.GenerateFullAnalysis(stockData)
// Returns: Executive, PriceAnalysis, Fundamentals, Risks, Outlook

// Or generate individual sections
executive, _ := aiClient.GenerateExecutiveSummary(stockData)
priceAnalysis, _ := aiClient.AnalyzePriceTrend(daily)
fundamentals, _ := aiClient.AnalyzeFundamentals(stockData)
risks, _ := aiClient.AssessRisks(stockData)
outlook, _ := aiClient.GenerateOutlook(stockData)

// Custom analysis
custom, _ := aiClient.CustomAnalysis(stockData, "What are the key growth drivers?")
```

### Available AI Models (OpenRouter)

Free models:
- `nvidia/nemotron-3-nano-30b-a3b:free`
- `meta-llama/llama-3.2-3b-instruct:free`

Paid models:
- `openai/gpt-4o-mini`
- `openai/gpt-4o`
- `anthropic/claude-3-haiku`
- `anthropic/claude-3-sonnet`

## PDF Report with AI Summary

```go
report := alphavintage.NewReportBuilder(alphavintage.DefaultReportOptions())
report.AddPageNumbers()

// Cover page
report.AddPage()
report.AddTitle("IBM Stock Analysis")
report.AddTimestamp()

// AI Summary page
report.AddPage()
report.AddAISummary(summary) // Adds all AI sections

// Or add individual AI insights
report.AddAIExecutiveSummary(summary.Executive)
report.AddAIInsight("Custom Analysis", customAnalysis)

// Charts and data
report.AddPage()
report.AddHeading("Price Analysis")
report.AddDailyPriceChart(daily, chartOpts)

report.Save("IBM_report.pdf")
```

## Chart Functions

```go
opts := alphavintage.ChartOptions{
    Width:      1000,
    Height:     500,
    Title:      "Chart Title",
    ShowVolume: true,
}

// Generate PNG files
alphavintage.GenerateDailyPriceChartToFile(daily, "price.png", opts)
alphavintage.GenerateCandlestickChartToFile(daily, "candle.png", opts)
alphavintage.GenerateEarningsChartToFile(earnings, "eps.png", opts)
alphavintage.GenerateCashFlowChartToFile(cashflow, "cashflow.png", opts)

// Compare multiple stocks
datasets := map[string]*alphavintage.TimeSeriesDailyResponse{
    "AAPL": appleDaily,
    "MSFT": msftDaily,
}
alphavintage.GenerateComparisonChartToFile(datasets, "compare.png", opts)
```

## Adding a Logo

```go
// Option 1: Set logo in ReportOptions (appears on all pages)
opts := alphavintage.DefaultReportOptions()
opts.LogoPath = "company_logo.png"
opts.LogoPosition = alphavintage.LogoTopRight  // LogoTopLeft, LogoTopRight, LogoTopCenter
opts.LogoWidthMM = 25
report := alphavintage.NewReportBuilder(opts)

// Option 2: Set logo after creating report
report := alphavintage.NewReportBuilder(alphavintage.DefaultReportOptions())
report.SetLogo("logo.png", alphavintage.LogoTopRight, 30)

// Option 3: Add logo to specific page only
report.AddPage()
report.AddLogo("logo.png", alphavintage.LogoTopLeft, 25)
```

## PDF Report Builder

```go
report := alphavintage.NewReportBuilder(opts)

// Text
report.AddTitle("Title")
report.AddSubtitle("Subtitle")
report.AddHeading("Section")
report.AddText("Paragraph")
report.AddBoldText("Bold")
report.AddItalicText("Italic")
report.AddBulletPoint("Item")
report.AddKeyValue("Label", "Value")

// Layout
report.AddPage()
report.AddLineBreak(10)
report.AddHorizontalLine()
report.AddPageNumbers()
report.AddTimestamp()

// Data
report.AddTable(headers, rows)
report.AddBalanceSheetSummary(balance)
report.AddCashFlowSummary(cashflow)
report.AddEarningsSummary(earnings, 5)
report.AddMarketStatusSummary(market)

// Charts
report.AddDailyPriceChart(daily, opts)
report.AddCandlestickChart(daily, opts)
report.AddEarningsChart(earnings, opts)
report.AddCashFlowChart(cashflow, opts)

// AI
report.AddAISummary(aiSummary)
report.AddAIInsight("Title", content)

// Save
report.Save("report.pdf")
```

## Environment Variables

```bash
export ALPHA_VANTAGE_API_KEY="your_key"
export OPENROUTER_API_KEY="your_key"  # For AI features
```

## Complete Example

```go
package main

import (
    "github.com/SwanHtetAungPhyo/alphavintage"
    "os"
    "time"
)

func main() {
    client := alphavintage.NewClient(os.Getenv("ALPHA_VANTAGE_API_KEY"))
    
    // Fetch data
    daily, _ := client.GetTimeSeriesDaily("AAPL", alphavintage.OutputSizeCompact)
    time.Sleep(12 * time.Second) // Rate limit
    earnings, _ := client.GetEarnings("AAPL")
    time.Sleep(12 * time.Second)
    cashflow, _ := client.GetCashFlow("AAPL")
    time.Sleep(12 * time.Second)
    balance, _ := client.GetBalanceSheet("AAPL")
    
    // AI Analysis
    ai := alphavintage.NewAIClient(alphavintage.AIConfig{
        APIKey: os.Getenv("OPENROUTER_API_KEY"),
        Model:  "nvidia/nemotron-3-nano-30b-a3b:free",
    })
    
    summary, _ := ai.GenerateFullAnalysis(alphavintage.StockAnalysisData{
        Symbol: "AAPL", Daily: daily, Earnings: earnings,
        CashFlow: cashflow, BalanceSheet: balance,
    })
    
    // Build PDF
    report := alphavintage.NewReportBuilder(alphavintage.DefaultReportOptions())
    report.AddPageNumbers()
    
    report.AddPage().AddTitle("AAPL Analysis").AddTimestamp()
    report.AddPage().AddAISummary(summary)
    report.AddPage().AddHeading("Price").AddDailyPriceChart(daily, alphavintage.ChartOptions{})
    report.AddPage().AddHeading("Earnings").AddEarningsChart(earnings, alphavintage.ChartOptions{})
    
    report.Save("AAPL_report.pdf")
}
```

## Financial Datasets API

The library also supports the Financial Datasets API for more comprehensive data:

```go
fd := alphavintage.NewFinancialDatasetsClient("YOUR_FD_API_KEY")

// Company info
company, _ := fd.GetCompanyFacts("AAPL")

// Real-time price
snapshot, _ := fd.GetPriceSnapshot("AAPL")

// Historical prices
prices, _ := fd.GetPrices("AAPL", alphavintage.FDIntervalDay, 1, "2024-01-01", "2024-12-31", 0)

// Financial statements
income, _ := fd.GetIncomeStatements("AAPL", alphavintage.FDPeriodAnnual, 5)
balance, _ := fd.GetBalanceSheets("AAPL", alphavintage.FDPeriodAnnual, 5)
cashflow, _ := fd.GetCashFlowStatements("AAPL", alphavintage.FDPeriodAnnual, 5)

// Financial metrics/ratios
metrics, _ := fd.GetFinancialMetricsSnapshot("AAPL")

// Insider trades
insiders, _ := fd.GetInsiderTrades("AAPL", 10)

// Institutional ownership
institutions, _ := fd.GetInstitutionalOwnership("AAPL", 10)

// News
news, _ := fd.GetNews("AAPL", "2024-01-01", "2024-12-31", 10)
```

### Financial Datasets PDF Report

```go
report := alphavintage.NewReportBuilder(opts)

// Company info
report.AddFDCompanyInfo(company)
report.AddFDPriceSnapshot(snapshot)

// Charts
report.AddFDPriceChart(prices, chartOpts)
report.AddFDRevenueChart(income, chartOpts)

// Financial data
report.AddFDIncomeStatementSummary(income, 5)
report.AddFDBalanceSheetSummary(balance)
report.AddFDCashFlowSummary(cashflow)
report.AddFDFinancialMetrics(metrics)

// Trading activity
report.AddFDInsiderTrades(insiders, 10)
report.AddFDInstitutionalOwnership(institutions, 10)
report.AddFDNews(news, 5)

report.Save("report.pdf")
```

## License

MIT
