package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SwanHtetAungPhyo/alphavintage"
)

func main() {
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		apiKey = "UPS6QRH073V81U5Z" // fallback
	}
	client := alphavintage.NewClient(apiKey)

	// Get market status
	fmt.Println("=== Market Status ===")
	status, err := client.GetMarketStatus()
	if err != nil {
		log.Printf("Market status error: %v", err)
	} else {
		for _, m := range status.Markets {
			fmt.Printf("%s (%s): %s\n", m.Region, m.MarketType, m.CurrentStatus)
		}
	}

	time.Sleep(12 * time.Second) // Rate limit: 5 calls/min for free tier

	// Get daily time series
	fmt.Println("\n=== Daily Time Series (IBM) ===")
	daily, err := client.GetTimeSeriesDaily("IBM", alphavintage.OutputSizeCompact)
	if err != nil {
		log.Printf("Daily series error: %v", err)
	} else {
		fmt.Printf("Symbol: %s, Last Refreshed: %s\n", daily.MetaData.Symbol, daily.MetaData.LastRefreshed)
		count := 0
		for date, data := range daily.TimeSeries {
			if count >= 3 {
				break
			}
			fmt.Printf("  %s: Open=%s, Close=%s, Volume=%s\n", date, data.Open, data.Close, data.Volume)
			count++
		}
	}

	time.Sleep(12 * time.Second)

	// Get intraday time series
	fmt.Println("\n=== Intraday Time Series (IBM, 5min) ===")
	intraday, err := client.GetTimeSeriesIntraday("IBM", alphavintage.Interval5Min, alphavintage.OutputSizeCompact)
	if err != nil {
		log.Printf("Intraday series error: %v", err)
	} else {
		fmt.Printf("Symbol: %s, Interval: %s\n", intraday.MetaData.Symbol, intraday.MetaData.Interval)
		count := 0
		for ts, data := range intraday.TimeSeries {
			if count >= 3 {
				break
			}
			fmt.Printf("  %s: Open=%s, Close=%s\n", ts, data.Open, data.Close)
			count++
		}
	}

	time.Sleep(12 * time.Second)

	// Get balance sheet
	fmt.Println("\n=== Balance Sheet (IBM) ===")
	balance, err := client.GetBalanceSheet("IBM")
	if err != nil {
		log.Printf("Balance sheet error: %v", err)
	} else if len(balance.AnnualReports) > 0 {
		report := balance.AnnualReports[0]
		fmt.Printf("Fiscal Date: %s\n", report.FiscalDateEnding)
		fmt.Printf("Total Assets: %s\n", report.TotalAssets)
		fmt.Printf("Total Liabilities: %s\n", report.TotalLiabilities)
	}

	time.Sleep(12 * time.Second)

	// Get cash flow
	fmt.Println("\n=== Cash Flow (IBM) ===")
	cashflow, err := client.GetCashFlow("IBM")
	if err != nil {
		log.Printf("Cash flow error: %v", err)
	} else if len(cashflow.AnnualReports) > 0 {
		report := cashflow.AnnualReports[0]
		fmt.Printf("Fiscal Date: %s\n", report.FiscalDateEnding)
		fmt.Printf("Operating Cashflow: %s\n", report.OperatingCashflow)
		fmt.Printf("Net Income: %s\n", report.NetIncome)
	}

	time.Sleep(12 * time.Second)

	// Get earnings
	fmt.Println("\n=== Earnings (IBM) ===")
	earnings, err := client.GetEarnings("IBM")
	if err != nil {
		log.Printf("Earnings error: %v", err)
	} else if len(earnings.AnnualEarnings) > 0 {
		for i, e := range earnings.AnnualEarnings {
			if i >= 3 {
				break
			}
			fmt.Printf("  %s: EPS=%s\n", e.FiscalDateEnding, e.ReportedEPS)
		}
	}

	time.Sleep(12 * time.Second)

	// Get news sentiment
	fmt.Println("\n=== News Sentiment (AAPL) ===")
	news, err := client.GetNewsSentiment(&alphavintage.NewsSentimentOptions{
		Tickers: "AAPL",
	})
	if err != nil {
		log.Printf("News sentiment error: %v", err)
	} else if len(news.Feed) > 0 {
		for i, item := range news.Feed {
			if i >= 3 {
				break
			}
			title := item.Title
			if len(title) > 60 {
				title = title[:60]
			}
			fmt.Printf("  %s (Sentiment: %s)\n", title, item.OverallSentimentLabel)
		}
	}
}
