package alphavintage

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// TestIntraday tests single-day intraday functionality
// Note: Intraday is a PREMIUM endpoint on Alpha Vantage
func TestIntraday() {
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		apiKey = "demo"
	}

	client := NewClient(apiKey)
	symbol := "IBM"

	fmt.Printf("=== Testing Intraday Data for %s ===\n\n", symbol)

	// Fetch intraday data
	fmt.Println("Fetching intraday data (5min intervals)...")
	fmt.Println("Note: This requires Alpha Vantage PREMIUM subscription")
	intraday, err := client.GetTimeSeriesIntraday(symbol, Interval5Min, OutputSizeFull)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("\nIntraday is a premium feature. Testing daily data instead...")
		TestSingleDayFromDaily()
		return
	}
	fmt.Printf("Total data points: %d\n", len(intraday.TimeSeries))

	// Try to filter for recent dates
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	fmt.Printf("\nFiltering for %s...\n", today)
	filtered := FilterIntradayByDate(intraday, today)
	fmt.Printf("Points for today: %d\n", len(filtered.TimeSeries))

	if len(filtered.TimeSeries) == 0 {
		fmt.Printf("\nFiltering for %s...\n", yesterday)
		filtered = FilterIntradayByDate(intraday, yesterday)
		fmt.Printf("Points for yesterday: %d\n", len(filtered.TimeSeries))
	}

	// Get summary if we have data
	if len(filtered.TimeSeries) > 0 {
		summary, err := GetIntradaySummary(filtered)
		if err != nil {
			fmt.Printf("Summary error: %v\n", err)
			return
		}

		fmt.Println("\n=== Intraday Summary ===")
		fmt.Printf("Symbol:      %s\n", summary.Symbol)
		fmt.Printf("Date:        %s\n", summary.Date)
		fmt.Printf("Interval:    %s\n", summary.Interval)
		fmt.Printf("Open:        $%.2f\n", summary.Open)
		fmt.Printf("High:        $%.2f\n", summary.High)
		fmt.Printf("Low:         $%.2f\n", summary.Low)
		fmt.Printf("Close:       $%.2f\n", summary.Close)
		fmt.Printf("Total Vol:   %d\n", summary.TotalVol)
		fmt.Printf("Data Points: %d\n", summary.DataPoints)

		// Generate chart
		fmt.Println("\nGenerating intraday chart...")
		chartOpts := ChartOptions{
			Title:      fmt.Sprintf("%s Intraday - %s", symbol, summary.Date),
			Width:      1200,
			Height:     600,
			ShowVolume: true,
		}

		err = GenerateIntradayChartToFile(filtered, symbol+"_intraday_test.png", chartOpts)
		if err != nil {
			fmt.Printf("Chart error: %v\n", err)
		} else {
			fmt.Printf("Saved: %s_intraday_test.png\n", symbol)
		}

		// Generate PDF
		fmt.Println("\nGenerating PDF report...")
		opts := DefaultReportOptions()
		opts.Title = symbol + " Intraday Report"
		report := NewReportBuilder(opts)
		report.AddPageNumbers()

		report.AddPage()
		report.AddTitle(symbol + " Intraday Analysis")
		report.AddSubtitle(summary.Date)
		report.AddTimestamp()

		report.AddHeading("Trading Summary")
		report.AddIntradaySummary(summary)

		report.AddHeading("Price Chart")
		report.AddIntradayChart(filtered, chartOpts)

		err = report.Save(symbol + "_intraday_test.pdf")
		if err != nil {
			fmt.Printf("PDF error: %v\n", err)
		} else {
			fmt.Printf("Saved: %s_intraday_test.pdf\n", symbol)
		}
	} else {
		fmt.Println("\nNo intraday data available for recent dates.")
		fmt.Println("This is normal for weekends/holidays or with demo API key.")
	}

	fmt.Println("\n=== Test Complete ===")
}

// TestSingleDayFromDaily tests getting single day data from daily time series (FREE tier)
func TestSingleDayFromDaily() {
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		apiKey = "demo"
	}

	client := NewClient(apiKey)
	symbol := "IBM"

	fmt.Printf("\n=== Testing Daily Data with Date Range for %s ===\n\n", symbol)

	// Fetch daily data (one API call)
	fmt.Println("Fetching daily data...")
	daily, err := client.GetTimeSeriesDaily(symbol, OutputSizeCompact)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Total days fetched: %d\n", len(daily.TimeSeries))

	// Get date range info
	oldestDate := GetOldestDate(daily)
	mostRecent := GetMostRecentDate(daily)
	fmt.Printf("Data range: %s to %s\n", oldestDate, mostRecent)

	// ============================================
	// Test 1: Get single day data
	// ============================================
	fmt.Println("\n--- Test 1: Single Day ---")
	point, ok := GetDailyDataPoint(daily, mostRecent)
	if ok {
		open, _ := strconv.ParseFloat(point.Open, 64)
		high, _ := strconv.ParseFloat(point.High, 64)
		low, _ := strconv.ParseFloat(point.Low, 64)
		close, _ := strconv.ParseFloat(point.Close, 64)
		fmt.Printf("Date: %s | O: $%.2f H: $%.2f L: $%.2f C: $%.2f\n",
			mostRecent, open, high, low, close)
	}

	// ============================================
	// Test 2: Filter last 5 trading days
	// ============================================
	fmt.Println("\n--- Test 2: Last 5 Trading Days ---")
	last5 := FilterDailyLastNDays(daily, 5)
	fmt.Printf("Filtered to %d days\n", len(last5.TimeSeries))

	summary5, _ := GetDailyRangeSummary(last5)
	if summary5 != nil {
		fmt.Printf("Period: %s to %s\n", summary5.StartDate, summary5.EndDate)
		fmt.Printf("Open: $%.2f -> Close: $%.2f (%.2f%%)\n",
			summary5.PeriodOpen, summary5.PeriodClose, summary5.PriceChangePct)
	}

	// ============================================
	// Test 3: Filter last 30 trading days
	// ============================================
	fmt.Println("\n--- Test 3: Last 30 Trading Days ---")
	last30 := FilterDailyLastNDays(daily, 30)
	fmt.Printf("Filtered to %d days\n", len(last30.TimeSeries))

	summary30, _ := GetDailyRangeSummary(last30)
	if summary30 != nil {
		fmt.Printf("Period: %s to %s\n", summary30.StartDate, summary30.EndDate)
		fmt.Printf("High: $%.2f (%s) | Low: $%.2f (%s)\n",
			summary30.PeriodHigh, summary30.HighDate, summary30.PeriodLow, summary30.LowDate)
		fmt.Printf("Change: $%.2f (%.2f%%)\n", summary30.PriceChange, summary30.PriceChangePct)
	}

	// ============================================
	// Test 4: Filter by specific date range
	// ============================================
	fmt.Println("\n--- Test 4: Custom Date Range ---")
	// Example: Filter for December 2025
	rangeData := FilterDailyByDateRange(daily, "2025-12-01", "2025-12-15")
	fmt.Printf("Dec 1-15, 2025: %d trading days\n", len(rangeData.TimeSeries))

	summaryRange, _ := GetDailyRangeSummary(rangeData)
	if summaryRange != nil && summaryRange.TradingDays > 0 {
		fmt.Printf("Period: %s to %s\n", summaryRange.StartDate, summaryRange.EndDate)
		fmt.Printf("Change: $%.2f (%.2f%%)\n", summaryRange.PriceChange, summaryRange.PriceChangePct)
	}

	// ============================================
	// Generate PDF Report
	// ============================================
	fmt.Println("\n--- Generating PDF Report ---")
	opts := DefaultReportOptions()
	opts.Title = symbol + " Date Range Analysis"
	report := NewReportBuilder(opts)
	report.AddPageNumbers()

	// Cover page
	report.AddPage()
	report.AddTitle(symbol + " Analysis")
	report.AddSubtitle("Date Range Report")
	report.AddTimestamp()

	// Last 5 days summary
	report.AddPage()
	report.AddHeading("Last 5 Trading Days")
	report.AddDailyRangeSummary(summary5)
	report.AddDailyPriceChart(last5, ChartOptions{
		Title:      symbol + " - Last 5 Days",
		ShowVolume: true,
	})

	// Last 30 days summary
	report.AddPage()
	report.AddHeading("Last 30 Trading Days")
	report.AddDailyRangeSummary(summary30)
	report.AddDailyPriceChart(last30, ChartOptions{
		Title:      symbol + " - Last 30 Days",
		ShowVolume: true,
	})

	// Full data
	report.AddPage()
	report.AddHeading("Full Data Range")
	fullSummary, _ := GetDailyRangeSummary(daily)
	report.AddDailyRangeSummary(fullSummary)
	report.AddDailyPriceChart(daily, ChartOptions{
		Title:      symbol + " - All Available Data",
		ShowVolume: true,
	})

	err = report.Save(symbol + "_daterange_test.pdf")
	if err != nil {
		fmt.Printf("PDF error: %v\n", err)
	} else {
		fmt.Printf("\nSaved: %s_daterange_test.pdf\n", symbol)
	}

	fmt.Println("\n=== Test Complete ===")
}
