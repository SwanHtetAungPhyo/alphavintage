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
		apiKey = "UPS6QRH073V81U5Z"
	}
	client := alphavintage.NewClient(apiKey)

	symbol := "IBM"

	// Fetch all data
	fmt.Printf("Fetching data for %s...\n", symbol)

	daily, err := client.GetTimeSeriesDaily(symbol, alphavintage.OutputSizeCompact)
	if err != nil {
		log.Printf("Daily series error: %v", err)
	} else {
		fmt.Printf("✓ Daily data: %d points\n", len(daily.TimeSeries))
	}

	time.Sleep(12 * time.Second)

	earnings, err := client.GetEarnings(symbol)
	if err != nil {
		log.Printf("Earnings error: %v", err)
	} else {
		fmt.Printf("✓ Earnings: %d annual records\n", len(earnings.AnnualEarnings))
	}

	time.Sleep(12 * time.Second)

	cashflow, err := client.GetCashFlow(symbol)
	if err != nil {
		log.Printf("Cash flow error: %v", err)
	} else {
		fmt.Printf("✓ Cash flow: %d annual records\n", len(cashflow.AnnualReports))
	}

	time.Sleep(12 * time.Second)

	balance, err := client.GetBalanceSheet(symbol)
	if err != nil {
		log.Printf("Balance sheet error: %v", err)
	} else {
		fmt.Printf("✓ Balance sheet: %d annual records\n", len(balance.AnnualReports))
	}

	time.Sleep(12 * time.Second)

	market, err := client.GetMarketStatus()
	if err != nil {
		log.Printf("Market status error: %v", err)
	} else {
		fmt.Printf("✓ Market status: %d markets\n", len(market.Markets))
	}

	// Create PDF Report
	fmt.Println("\nGenerating PDF report...")

	opts := alphavintage.DefaultReportOptions()
	opts.Title = fmt.Sprintf("%s Stock Analysis Report", symbol)
	opts.Author = "Alpha Vantage Go Client"

	report := alphavintage.NewReportBuilder(opts)
	report.AddPageNumbers()

	// Cover page
	report.AddPage()
	report.AddLineBreak(40)
	report.AddTitle(fmt.Sprintf("%s Stock Analysis", symbol))
	report.AddLineBreak(10)
	report.AddSubtitle("Comprehensive Financial Report")
	report.AddLineBreak(20)
	report.AddTimestamp()
	report.AddLineBreak(10)
	report.AddText("This report provides a comprehensive analysis of the stock including price trends, earnings history, cash flow analysis, and balance sheet summary.")

	// Price Analysis Page
	report.AddPage()
	report.AddHeading("Price Analysis")
	report.AddText(fmt.Sprintf("The following chart shows the daily closing prices for %s over the recent trading period.", symbol))
	report.AddLineBreak(5)

	chartOpts := alphavintage.ChartOptions{
		Title:      fmt.Sprintf("%s Daily Price", symbol),
		Width:      800,
		Height:     400,
		ShowVolume: true,
	}
	report.AddDailyPriceChart(daily, chartOpts)

	report.AddLineBreak(5)
	report.AddHeading("Price Range Analysis")
	chartOpts.Title = fmt.Sprintf("%s High/Low/Close", symbol)
	report.AddCandlestickChart(daily, chartOpts)

	// Earnings Page
	report.AddPage()
	report.AddHeading("Earnings Analysis")
	report.AddText("Annual earnings per share (EPS) provides insight into the company's profitability over time.")
	report.AddLineBreak(5)

	chartOpts.Title = fmt.Sprintf("%s Annual EPS", symbol)
	chartOpts.Width = 700
	chartOpts.Height = 350
	report.AddEarningsChart(earnings, chartOpts)

	report.AddLineBreak(5)
	report.AddEarningsSummary(earnings, 5)

	// Cash Flow Page
	report.AddPage()
	report.AddHeading("Cash Flow Analysis")
	report.AddText("Cash flow analysis shows how the company generates and uses cash across operating, investing, and financing activities.")
	report.AddLineBreak(5)

	chartOpts.Title = fmt.Sprintf("%s Cash Flow Trends", symbol)
	report.AddCashFlowChart(cashflow, chartOpts)

	report.AddLineBreak(5)
	report.AddCashFlowSummary(cashflow)

	// Balance Sheet Page
	report.AddPage()
	report.AddHeading("Balance Sheet Summary")
	report.AddText("The balance sheet provides a snapshot of the company's financial position at a specific point in time.")
	report.AddLineBreak(5)
	report.AddBalanceSheetSummary(balance)

	// Market Status Page
	report.AddPage()
	report.AddHeading("Global Market Status")
	report.AddText("Current status of major global markets.")
	report.AddLineBreak(5)
	report.AddMarketStatusSummary(market)

	// Disclaimer Page
	report.AddPage()
	report.AddHeading("Disclaimer")
	report.AddItalicText("This report is generated automatically using data from Alpha Vantage API. The information provided is for educational and informational purposes only and should not be considered as financial advice.")
	report.AddLineBreak(5)
	report.AddText("Past performance is not indicative of future results. Always conduct your own research and consult with a qualified financial advisor before making investment decisions.")

	// Save the report
	filename := fmt.Sprintf("%s_report.pdf", symbol)
	err = report.Save(filename)
	if err != nil {
		log.Fatalf("Failed to save PDF: %v", err)
	}

	fmt.Printf("\n✓ Report saved: %s\n", filename)
}
