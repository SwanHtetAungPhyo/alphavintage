package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SwanHtetAungPhyo/alphavintage"
)

func main() {
	// API Keys
	avKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if avKey == "" {
		avKey = "UPS6QRH073V81U5Z"
	}
	orKey := os.Getenv("OPENROUTER_API_KEY")
	if orKey == "" {
		orKey = "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"
	}
	fdKey := "4959d8c9-f9c4-4633-bb02-4c98d0f8923f"

	symbol := "AAPL"
	fmt.Printf("=== Alpha Vantage Go Library - Full Test ===\n\n")

	// ==========================================
	// PART 1: Alpha Vantage API
	// ==========================================
	fmt.Println("--- PART 1: Alpha Vantage API ---")
	client := alphavintage.NewClient(avKey)

	// Daily prices
	fmt.Printf("Fetching daily prices for %s...\n", symbol)
	daily, err := client.GetTimeSeriesDaily(symbol, alphavintage.OutputSizeCompact)
	if err != nil {
		log.Printf("Daily error: %v", err)
	} else {
		fmt.Printf("  Daily: %d data points\n", len(daily.TimeSeries))
	}
	time.Sleep(12 * time.Second)

	// Earnings
	fmt.Println("Fetching earnings...")
	earnings, err := client.GetEarnings(symbol)
	if err != nil {
		log.Printf("Earnings error: %v", err)
	} else {
		fmt.Printf("  Annual earnings: %d records\n", len(earnings.AnnualEarnings))
		fmt.Printf("  Quarterly earnings: %d records\n", len(earnings.QuarterlyEarnings))
	}
	time.Sleep(12 * time.Second)

	// Cash Flow
	fmt.Println("Fetching cash flow...")
	cashflow, err := client.GetCashFlow(symbol)
	if err != nil {
		log.Printf("Cash flow error: %v", err)
	} else {
		fmt.Printf("  Annual reports: %d\n", len(cashflow.AnnualReports))
	}
	time.Sleep(12 * time.Second)

	// Balance Sheet
	fmt.Println("Fetching balance sheet...")
	balance, err := client.GetBalanceSheet(symbol)
	if err != nil {
		log.Printf("Balance sheet error: %v", err)
	} else {
		fmt.Printf("  Annual reports: %d\n", len(balance.AnnualReports))
	}
	time.Sleep(12 * time.Second)

	// Market Status
	fmt.Println("Fetching market status...")
	market, err := client.GetMarketStatus()
	if err != nil {
		log.Printf("Market status error: %v", err)
	} else {
		fmt.Printf("  Markets: %d\n", len(market.Markets))
	}
	time.Sleep(12 * time.Second)

	// News Sentiment
	fmt.Println("Fetching news sentiment...")
	newsOpts := &alphavintage.NewsSentimentOptions{Tickers: symbol, Limit: 5}
	news, err := client.GetNewsSentiment(newsOpts)
	if err != nil {
		log.Printf("News error: %v", err)
	} else {
		fmt.Printf("  News articles: %d\n", len(news.Feed))
	}

	// ==========================================
	// PART 2: Chart Generation (PNG files)
	// ==========================================
	fmt.Println("\n--- PART 2: Chart Generation ---")

	chartOpts := alphavintage.DefaultChartOptions()

	// Daily Price Chart
	if daily != nil && len(daily.TimeSeries) > 0 {
		chartOpts.Title = symbol + " Daily Price"
		err := alphavintage.GenerateDailyPriceChartToFile(daily, symbol+"_price.png", chartOpts)
		if err != nil {
			log.Printf("Price chart error: %v", err)
		} else {
			fmt.Println("  Generated: " + symbol + "_price.png")
		}

		// Candlestick Chart
		chartOpts.Title = symbol + " Candlestick"
		err = alphavintage.GenerateCandlestickChartToFile(daily, symbol+"_candlestick.png", chartOpts)
		if err != nil {
			log.Printf("Candlestick error: %v", err)
		} else {
			fmt.Println("  Generated: " + symbol + "_candlestick.png")
		}
	}

	// Earnings Chart
	if earnings != nil && len(earnings.AnnualEarnings) > 0 {
		chartOpts.Title = symbol + " Annual EPS"
		err := alphavintage.GenerateEarningsChartToFile(earnings, symbol+"_earnings.png", chartOpts)
		if err != nil {
			log.Printf("Earnings chart error: %v", err)
		} else {
			fmt.Println("  Generated: " + symbol + "_earnings.png")
		}
	}

	// Cash Flow Chart
	if cashflow != nil && len(cashflow.AnnualReports) > 0 {
		chartOpts.Title = symbol + " Cash Flow"
		err := alphavintage.GenerateCashFlowChartToFile(cashflow, symbol+"_cashflow.png", chartOpts)
		if err != nil {
			log.Printf("Cash flow chart error: %v", err)
		} else {
			fmt.Println("  Generated: " + symbol + "_cashflow.png")
		}
	}

	// ==========================================
	// PART 3: AI Analysis
	// ==========================================
	fmt.Println("\n--- PART 3: AI Analysis ---")
	var aiSummary *alphavintage.AnalysisSummary
	var customInsight string

	if orKey != "" {
		ai := alphavintage.NewAIClient(alphavintage.AIConfig{
			APIKey:    orKey,
			Model:     "nvidia/nemotron-3-nano-30b-a3b:free",
			Reasoning: false,
		})

		stockData := alphavintage.StockAnalysisData{
			Symbol:       symbol,
			Daily:        daily,
			Earnings:     earnings,
			CashFlow:     cashflow,
			BalanceSheet: balance,
			News:         news,
		}

		fmt.Println("Generating full AI analysis...")
		aiSummary, err = ai.GenerateFullAnalysis(stockData)
		if err != nil {
			log.Printf("AI analysis error: %v", err)
		} else {
			fmt.Println("  Executive summary: generated")
			fmt.Println("  Price analysis: generated")
			fmt.Println("  Fundamentals: generated")
			fmt.Println("  Risk assessment: generated")
			fmt.Println("  Outlook: generated")
		}

		// Custom analysis
		fmt.Println("Generating custom insight...")
		customInsight, err = ai.CustomAnalysis(stockData, "What are the top 3 investment considerations for this stock?")
		if err != nil {
			log.Printf("Custom analysis error: %v", err)
		} else {
			fmt.Println("  Custom insight: generated")
		}

		// News summary
		if news != nil && len(news.Feed) > 0 {
			fmt.Println("Summarizing news...")
			newsSummary, err := ai.SummarizeNews(news)
			if err != nil {
				log.Printf("News summary error: %v", err)
			} else {
				fmt.Printf("  News summary: %s...\n", truncate(newsSummary, 50))
			}
		}
	} else {
		fmt.Println("  Skipped (no OPENROUTER_API_KEY)")
	}

	// ==========================================
	// PART 4: Financial Datasets API
	// ==========================================
	fmt.Println("\n--- PART 4: Financial Datasets API ---")

	var fdCompany *alphavintage.FDCompanyFacts
	var fdSnapshot *alphavintage.FDPriceSnapshot
	var fdPrices []alphavintage.FDPrice
	var fdIncome []alphavintage.FDIncomeStatement
	var fdBalance []alphavintage.FDBalanceSheet
	var fdCashflow []alphavintage.FDCashFlowStatement
	var fdMetrics *alphavintage.FDFinancialMetrics
	var fdInsiders []alphavintage.FDInsiderTrade
	var fdInstitutions []alphavintage.FDInstitutionalOwnership
	var fdNews []alphavintage.FDNews

	if fdKey != "" {
		fd := alphavintage.NewFinancialDatasetsClient(fdKey)

		// Company Facts
		fmt.Println("Fetching company facts...")
		fdCompany, err = fd.GetCompanyFacts(symbol)
		if err != nil {
			log.Printf("FD company error: %v", err)
		} else {
			fmt.Printf("  Company: %s (%s)\n", fdCompany.Name, fdCompany.Sector)
		}

		// Price Snapshot
		fmt.Println("Fetching price snapshot...")
		fdSnapshot, err = fd.GetPriceSnapshot(symbol)
		if err != nil {
			log.Printf("FD snapshot error: %v", err)
		} else {
			fmt.Printf("  Price: $%.2f (%.2f%%)\n", fdSnapshot.Price, fdSnapshot.DayChangePercent)
		}

		// Historical Prices
		fmt.Println("Fetching historical prices...")
		endDate := time.Now().Format("2006-01-02")
		startDate := time.Now().AddDate(0, -3, 0).Format("2006-01-02")
		fdPrices, err = fd.GetPrices(symbol, alphavintage.FDIntervalDay, 1, startDate, endDate, 100)
		if err != nil {
			log.Printf("FD prices error: %v", err)
		} else {
			fmt.Printf("  Prices: %d data points\n", len(fdPrices))
		}

		// Income Statements
		fmt.Println("Fetching income statements...")
		fdIncome, err = fd.GetIncomeStatements(symbol, alphavintage.FDPeriodAnnual, 5)
		if err != nil {
			log.Printf("FD income error: %v", err)
		} else {
			fmt.Printf("  Income statements: %d\n", len(fdIncome))
		}

		// Balance Sheets
		fmt.Println("Fetching balance sheets...")
		fdBalance, err = fd.GetBalanceSheets(symbol, alphavintage.FDPeriodAnnual, 5)
		if err != nil {
			log.Printf("FD balance error: %v", err)
		} else {
			fmt.Printf("  Balance sheets: %d\n", len(fdBalance))
		}

		// Cash Flow Statements
		fmt.Println("Fetching cash flow statements...")
		fdCashflow, err = fd.GetCashFlowStatements(symbol, alphavintage.FDPeriodAnnual, 5)
		if err != nil {
			log.Printf("FD cashflow error: %v", err)
		} else {
			fmt.Printf("  Cash flow statements: %d\n", len(fdCashflow))
		}

		// Financial Metrics
		fmt.Println("Fetching financial metrics...")
		fdMetrics, err = fd.GetFinancialMetricsSnapshot(symbol)
		if err != nil {
			log.Printf("FD metrics error: %v", err)
		} else {
			fmt.Printf("  P/E Ratio: %.2f\n", fdMetrics.PriceToEarningsRatio)
		}

		// Insider Trades
		fmt.Println("Fetching insider trades...")
		fdInsiders, err = fd.GetInsiderTrades(symbol, 10)
		if err != nil {
			log.Printf("FD insiders error: %v", err)
		} else {
			fmt.Printf("  Insider trades: %d\n", len(fdInsiders))
		}

		// Institutional Ownership
		fmt.Println("Fetching institutional ownership...")
		fdInstitutions, err = fd.GetInstitutionalOwnership(symbol, 10)
		if err != nil {
			log.Printf("FD institutions error: %v", err)
		} else {
			fmt.Printf("  Institutional holders: %d\n", len(fdInstitutions))
		}

		// News
		fmt.Println("Fetching FD news...")
		fdNews, err = fd.GetNews(symbol, startDate, endDate, 5)
		if err != nil {
			log.Printf("FD news error: %v", err)
		} else {
			fmt.Printf("  News articles: %d\n", len(fdNews))
		}

		// FD Charts
		if len(fdPrices) > 0 {
			chartOpts.Title = symbol + " FD Price"
			err = alphavintage.GenerateFDPriceChartToFile(fdPrices, symbol+"_fd_price.png", chartOpts)
			if err != nil {
				log.Printf("FD price chart error: %v", err)
			} else {
				fmt.Println("  Generated: " + symbol + "_fd_price.png")
			}
		}

		if len(fdIncome) > 0 {
			chartOpts.Title = symbol + " Revenue"
			err = alphavintage.GenerateFDRevenueChartToFile(fdIncome, symbol+"_fd_revenue.png", chartOpts)
			if err != nil {
				log.Printf("FD revenue chart error: %v", err)
			} else {
				fmt.Println("  Generated: " + symbol + "_fd_revenue.png")
			}
		}
	} else {
		fmt.Println("  Skipped (no FINANCIAL_DATASETS_API_KEY)")
	}

	// ==========================================
	// PART 5: PDF Report Generation
	// ==========================================
	fmt.Println("\n--- PART 5: PDF Report Generation ---")

	opts := alphavintage.DefaultReportOptions()
	opts.Title = symbol + " Comprehensive Stock Analysis"
	opts.Author = "Alpha Vantage Go Library"
	opts.LogoPath = "logo.jpeg"
	opts.LogoPosition = alphavintage.LogoTopRight
	opts.LogoWidthMM = 25

	report := alphavintage.NewReportBuilder(opts)
	report.AddPageNumbers()

	// --- Cover Page ---
	report.AddPage()
	report.AddLineBreak(40)
	report.AddTitle(symbol)
	report.AddSubtitle("Comprehensive Stock Analysis Report")
	report.AddLineBreak(20)
	report.AddTimestamp()
	report.AddLineBreak(15)
	report.AddText("This report includes price analysis, fundamental data, AI-powered insights, and market intelligence from multiple data sources.")
	report.AddLineBreak(10)
	report.AddHorizontalLine()
	report.AddLineBreak(5)
	report.AddBoldText("Data Sources:")
	report.AddBulletPoint("Alpha Vantage API - Market data, fundamentals")
	report.AddBulletPoint("Financial Datasets API - Enhanced financial data")
	report.AddBulletPoint("OpenRouter AI - Intelligent analysis")

	// --- Table of Contents ---
	report.AddPage()
	report.AddHeading("Table of Contents")
	report.AddLineBreak(5)
	report.AddNumberedItem(1, "AI Executive Summary")
	report.AddNumberedItem(2, "Price Analysis")
	report.AddNumberedItem(3, "Earnings Analysis")
	report.AddNumberedItem(4, "Cash Flow Analysis")
	report.AddNumberedItem(5, "Balance Sheet")
	report.AddNumberedItem(6, "Financial Metrics")
	report.AddNumberedItem(7, "Insider Trading Activity")
	report.AddNumberedItem(8, "Institutional Ownership")
	report.AddNumberedItem(9, "News & Sentiment")
	report.AddNumberedItem(10, "Global Market Status")
	report.AddNumberedItem(11, "Disclaimer")

	// --- AI Summary ---
	if aiSummary != nil {
		report.AddPage()
		report.AddTitle("AI Analysis Summary")
		report.AddLineBreak(5)
		report.AddAISummary(aiSummary)

		if customInsight != "" {
			report.AddPage()
			report.AddAIInsight("Investment Considerations", customInsight)
		}
	}

	// --- Price Analysis ---
	report.AddPage()
	report.AddHeading("1. Price Analysis")
	report.AddText("Historical daily price data with volume analysis.")
	if daily != nil && len(daily.TimeSeries) > 0 {
		chartOpts.Title = symbol + " Daily Price"
		chartOpts.ShowVolume = true
		report.AddDailyPriceChart(daily, chartOpts)

		report.AddLineBreak(10)
		report.AddHeading("Price Range (High/Low/Close)")
		chartOpts.Title = symbol + " H/L/C"
		report.AddCandlestickChart(daily, chartOpts)
	} else {
		report.AddItalicText("Price data not available.")
	}

	// --- Earnings ---
	report.AddPage()
	report.AddHeading("2. Earnings Analysis")
	if earnings != nil && len(earnings.AnnualEarnings) > 0 {
		chartOpts.Title = symbol + " Annual EPS"
		report.AddEarningsChart(earnings, chartOpts)
		report.AddLineBreak(10)
		report.AddBoldText("Annual Earnings History")
		report.AddEarningsSummary(earnings, 5)
	} else {
		report.AddItalicText("Earnings data not available.")
	}

	// --- Cash Flow ---
	report.AddPage()
	report.AddHeading("3. Cash Flow Analysis")
	if cashflow != nil && len(cashflow.AnnualReports) > 0 {
		chartOpts.Title = symbol + " Cash Flow Trends"
		report.AddCashFlowChart(cashflow, chartOpts)
		report.AddLineBreak(10)
		report.AddBoldText("Latest Cash Flow Summary")
		report.AddCashFlowSummary(cashflow)
	} else {
		report.AddItalicText("Cash flow data not available.")
	}

	// --- Balance Sheet ---
	report.AddPage()
	report.AddHeading("4. Balance Sheet")
	if balance != nil && len(balance.AnnualReports) > 0 {
		report.AddBalanceSheetSummary(balance)
	} else {
		report.AddItalicText("Balance sheet data not available.")
	}

	// --- Financial Datasets Data ---
	if fdKey != "" {
		// Company Info & Metrics
		report.AddPage()
		report.AddHeading("5. Company Profile & Metrics")
		if fdCompany != nil {
			report.AddFDCompanyInfo(fdCompany)
		}
		if fdSnapshot != nil {
			report.AddBoldText("Current Price")
			report.AddFDPriceSnapshot(fdSnapshot)
		}
		if fdMetrics != nil {
			report.AddBoldText("Financial Ratios")
			report.AddFDFinancialMetrics(fdMetrics)
		}

		// FD Price Chart
		if len(fdPrices) > 0 {
			report.AddPage()
			report.AddHeading("6. Recent Price History (FD)")
			chartOpts.Title = symbol + " 3-Month Price"
			report.AddFDPriceChart(fdPrices, chartOpts)
		}

		// Income Statement
		if len(fdIncome) > 0 {
			report.AddPage()
			report.AddHeading("7. Income Statement")
			chartOpts.Title = symbol + " Revenue Trend"
			report.AddFDRevenueChart(fdIncome, chartOpts)
			report.AddLineBreak(10)
			report.AddFDIncomeStatementSummary(fdIncome, 5)
		}

		// FD Balance Sheet
		if len(fdBalance) > 0 {
			report.AddPage()
			report.AddHeading("8. Balance Sheet (FD)")
			report.AddFDBalanceSheetSummary(fdBalance)
		}

		// FD Cash Flow
		if len(fdCashflow) > 0 {
			report.AddPage()
			report.AddHeading("9. Cash Flow (FD)")
			report.AddFDCashFlowSummary(fdCashflow)
		}

		// Insider Trades
		if len(fdInsiders) > 0 {
			report.AddPage()
			report.AddHeading("10. Insider Trading Activity")
			report.AddText("Recent insider transactions:")
			report.AddFDInsiderTrades(fdInsiders, 10)
		}

		// Institutional Ownership
		if len(fdInstitutions) > 0 {
			report.AddPage()
			report.AddHeading("11. Institutional Ownership")
			report.AddText("Top institutional holders:")
			report.AddFDInstitutionalOwnership(fdInstitutions, 10)
		}

		// FD News
		if len(fdNews) > 0 {
			report.AddPage()
			report.AddHeading("12. Recent News (FD)")
			report.AddFDNews(fdNews, 5)
		}
	}

	// --- Market Status ---
	report.AddPage()
	report.AddHeading("Global Market Status")
	if market != nil && len(market.Markets) > 0 {
		report.AddMarketStatusSummary(market)
	} else {
		report.AddItalicText("Market status not available.")
	}

	// --- Custom Table Example ---
	report.AddPage()
	report.AddHeading("Summary Table Example")
	report.AddText("Custom table demonstration:")
	headers := []string{"Metric", "Value", "Status"}
	rows := [][]string{
		{"Daily Data Points", fmt.Sprintf("%d", len(daily.TimeSeries)), "OK"},
		{"Earnings Records", fmt.Sprintf("%d", len(earnings.AnnualEarnings)), "OK"},
		{"Cash Flow Reports", fmt.Sprintf("%d", len(cashflow.AnnualReports)), "OK"},
		{"Balance Sheets", fmt.Sprintf("%d", len(balance.AnnualReports)), "OK"},
	}
	report.AddTable(headers, rows)

	// --- Disclaimer ---
	report.AddPage()
	report.AddHeading("Disclaimer")
	report.AddLineBreak(5)
	report.AddItalicText("IMPORTANT: This report is generated for educational and informational purposes only. It does not constitute financial advice, investment recommendations, or an offer to buy or sell any securities.")
	report.AddLineBreak(5)
	report.AddText("The data presented in this report is sourced from third-party APIs (Alpha Vantage, Financial Datasets) and may contain errors or delays. AI-generated content is produced by language models and may contain inaccuracies.")
	report.AddLineBreak(5)
	report.AddText("Always conduct your own research and consult with a qualified financial advisor before making investment decisions.")
	report.AddLineBreak(10)
	report.AddHorizontalLine()
	report.AddLineBreak(5)
	report.AddKeyValue("Report Generated", time.Now().Format("January 2, 2006 at 3:04 PM"))
	report.AddKeyValue("Library Version", "1.0.0")
	report.AddKeyValue("Data Sources", "Alpha Vantage, Financial Datasets, OpenRouter")

	// Test SaveToBytes (must be called BEFORE Save, as Save closes the PDF)
	pdfBytes, err := report.SaveToBytes()
	if err != nil {
		log.Printf("SaveToBytes error: %v", err)
	} else {
		fmt.Printf("  PDF size: %d bytes (ready for S3 upload)\n", len(pdfBytes))
	}

	// Save PDF to file
	filename := symbol + "_full_report.pdf"
	if err := report.Save(filename); err != nil {
		log.Fatalf("Failed to save PDF: %v", err)
	}
	fmt.Printf("  Saved: %s\n", filename)

	fmt.Println("\n=== Test Complete ===")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
