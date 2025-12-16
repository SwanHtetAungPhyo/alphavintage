package main

import (
	"fmt"
	"os"
	"time"

	"github.com/SwanHtetAungPhyo/alphavintage"
)

// IntradayExample demonstrates single-day intraday data analysis
// Note: Alpha Vantage FREE tier only provides last 1-2 trading days
// Premium subscription required for extended intraday history

func IntradayExample() {
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		apiKey = "ADDDFCJFW22S60KW" // Demo key has limited data
	}

	client := alphavintage.NewClient(apiKey)
	symbol := "IBM"

	fmt.Printf("Fetching intraday data for %s...\n", symbol)

	// Method 1: Get all intraday data and filter by date
	intraday, err := client.GetTimeSeriesIntraday(symbol, alphavintage.Interval5Min, alphavintage.OutputSizeFull)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Total intraday points: %d\n", len(intraday.TimeSeries))

	// Filter for today's date (or most recent trading day)
	today := time.Now().Format("2006-01-02")
	filtered := alphavintage.FilterIntradayByDate(intraday, today)
	fmt.Printf("Points for %s: %d\n", today, len(filtered.TimeSeries))

	// If no data for today, try yesterday
	if len(filtered.TimeSeries) == 0 {
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		filtered = alphavintage.FilterIntradayByDate(intraday, yesterday)
		fmt.Printf("Points for %s: %d\n", yesterday, len(filtered.TimeSeries))
	}

	// Method 2: Use GetSingleDayData helper (combines fetch + filter)
	// singleDay, err := client.GetSingleDayData(symbol, "2024-12-16", alphavintage.Interval5Min)

	// Get summary statistics
	if len(filtered.TimeSeries) > 0 {
		summary, err := alphavintage.GetIntradaySummary(filtered)
		if err != nil {
			fmt.Printf("Summary error: %v\n", err)
			return
		}

		fmt.Println("\n=== Intraday Summary ===")
		fmt.Printf("Symbol:     %s\n", summary.Symbol)
		fmt.Printf("Date:       %s\n", summary.Date)
		fmt.Printf("Interval:   %s\n", summary.Interval)
		fmt.Printf("Open:       $%.2f\n", summary.Open)
		fmt.Printf("High:       $%.2f\n", summary.High)
		fmt.Printf("Low:        $%.2f\n", summary.Low)
		fmt.Printf("Close:      $%.2f\n", summary.Close)
		fmt.Printf("Volume:     %d\n", summary.TotalVol)
		fmt.Printf("Data Points: %d\n", summary.DataPoints)

		// Generate intraday chart
		chartOpts := alphavintage.ChartOptions{
			Title:      fmt.Sprintf("%s Intraday - %s", symbol, summary.Date),
			Width:      1200,
			Height:     600,
			ShowVolume: true,
		}

		err = alphavintage.GenerateIntradayChartToFile(filtered, symbol+"_intraday.png", chartOpts)
		if err != nil {
			fmt.Printf("Chart error: %v\n", err)
		} else {
			fmt.Printf("\nSaved: %s_intraday.png\n", symbol)
		}

		// Generate PDF report with intraday data
		opts := alphavintage.DefaultReportOptions()
		opts.Title = symbol + " Intraday Analysis"
		report := alphavintage.NewReportBuilder(opts)
		report.AddPageNumbers()

		report.AddPage()
		report.AddTitle(symbol + " Intraday Analysis")
		report.AddSubtitle(summary.Date)
		report.AddTimestamp()

		report.AddHeading("Trading Summary")
		report.AddIntradaySummary(summary)

		report.AddHeading("Price Chart")
		report.AddIntradayChart(filtered, chartOpts)

		err = report.Save(symbol + "_intraday_report.pdf")
		if err != nil {
			fmt.Printf("PDF error: %v\n", err)
		} else {
			fmt.Printf("Saved: %s_intraday_report.pdf\n", symbol)
		}
	}
}
