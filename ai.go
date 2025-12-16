package alphavintage

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const openRouterURL = "https://openrouter.ai/api/v1/chat/completions"

// AIClient handles AI-powered analysis using OpenRouter
type AIClient struct {
	apiKey    string
	model     string
	resty     *resty.Client
	reasoning bool
}

// AIConfig configures the AI client
type AIConfig struct {
	APIKey    string
	Model     string // e.g., "nvidia/nemotron-3-nano-30b-a3b:free", "openai/gpt-4o-mini"
	Reasoning bool   // Enable reasoning for supported models
}

// DefaultAIConfig returns default AI configuration
func DefaultAIConfig() AIConfig {
	return AIConfig{
		Model:     "nvidia/nemotron-3-nano-30b-a3b:free",
		Reasoning: false,
	}
}

// NewAIClient creates a new AI client for OpenRouter
func NewAIClient(config AIConfig) *AIClient {
	if config.Model == "" {
		config.Model = "nvidia/nemotron-3-nano-30b-a3b:free"
	}
	return &AIClient{
		apiKey:    config.APIKey,
		model:     config.Model,
		resty:     resty.New().SetTimeout(60 * time.Second),
		reasoning: config.Reasoning,
	}
}

// SetModel changes the AI model
func (ai *AIClient) SetModel(model string) *AIClient {
	ai.model = model
	return ai
}

// SetReasoning enables/disables reasoning
func (ai *AIClient) SetReasoning(enabled bool) *AIClient {
	ai.reasoning = enabled
	return ai
}

type openRouterRequest struct {
	Model     string          `json:"model"`
	Messages  []aiMessage     `json:"messages"`
	Reasoning *reasoningOpts  `json:"reasoning,omitempty"`
}

type aiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type reasoningOpts struct {
	Enabled bool `json:"enabled"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (ai *AIClient) chat(prompt string) (string, error) {
	req := openRouterRequest{
		Model: ai.model,
		Messages: []aiMessage{
			{Role: "user", Content: prompt},
		},
	}

	if ai.reasoning {
		req.Reasoning = &reasoningOpts{Enabled: true}
	}

	resp, err := ai.resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+ai.apiKey).
		SetBody(req).
		Post(openRouterURL)

	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	var result openRouterResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return result.Choices[0].Message.Content, nil
}


// StockAnalysisData holds all data for AI analysis
type StockAnalysisData struct {
	Symbol       string
	Daily        *TimeSeriesDailyResponse
	Earnings     *EarningsResponse
	CashFlow     *CashFlowResponse
	BalanceSheet *BalanceSheetResponse
	News         *NewsSentimentResponse
}

// AnalysisSummary contains AI-generated summaries
type AnalysisSummary struct {
	Executive    string // Executive summary
	PriceAnalysis string // Price trend analysis
	Fundamentals string // Fundamental analysis
	Risks        string // Risk assessment
	Outlook      string // Future outlook
}

// GenerateFullAnalysis generates comprehensive AI analysis
func (ai *AIClient) GenerateFullAnalysis(data StockAnalysisData) (*AnalysisSummary, error) {
	summary := &AnalysisSummary{}
	var err error

	// Generate each section
	summary.Executive, err = ai.GenerateExecutiveSummary(data)
	if err != nil {
		summary.Executive = "Unable to generate executive summary."
	}

	summary.PriceAnalysis, err = ai.AnalyzePriceTrend(data.Daily)
	if err != nil {
		summary.PriceAnalysis = "Unable to analyze price trends."
	}

	summary.Fundamentals, err = ai.AnalyzeFundamentals(data)
	if err != nil {
		summary.Fundamentals = "Unable to analyze fundamentals."
	}

	summary.Risks, err = ai.AssessRisks(data)
	if err != nil {
		summary.Risks = "Unable to assess risks."
	}

	summary.Outlook, err = ai.GenerateOutlook(data)
	if err != nil {
		summary.Outlook = "Unable to generate outlook."
	}

	return summary, nil
}

// GenerateExecutiveSummary creates a brief executive summary
func (ai *AIClient) GenerateExecutiveSummary(data StockAnalysisData) (string, error) {
	prompt := fmt.Sprintf(`Analyze this stock data for %s and provide a brief executive summary (3-4 sentences).

%s

Provide a concise, professional summary focusing on key metrics and overall health.`, 
		data.Symbol, formatDataForAI(data))

	return ai.chat(prompt)
}

// AnalyzePriceTrend analyzes price movements
func (ai *AIClient) AnalyzePriceTrend(data *TimeSeriesDailyResponse) (string, error) {
	if data == nil || len(data.TimeSeries) == 0 {
		return "", fmt.Errorf("no price data")
	}

	priceData := extractPriceSummary(data)
	prompt := fmt.Sprintf(`Analyze this stock price data and provide insights (3-4 sentences):

%s

Focus on: trend direction, volatility, support/resistance levels, and notable patterns.`, priceData)

	return ai.chat(prompt)
}

// AnalyzeFundamentals analyzes earnings, cash flow, balance sheet
func (ai *AIClient) AnalyzeFundamentals(data StockAnalysisData) (string, error) {
	fundamentals := formatFundamentalsForAI(data)
	prompt := fmt.Sprintf(`Analyze these fundamentals for %s (3-4 sentences):

%s

Focus on: profitability trends, financial health, and key ratios.`, data.Symbol, fundamentals)

	return ai.chat(prompt)
}

// AssessRisks identifies potential risks
func (ai *AIClient) AssessRisks(data StockAnalysisData) (string, error) {
	riskData := formatRiskDataForAI(data)
	prompt := fmt.Sprintf(`Identify key risks for %s based on this data (3-4 bullet points):

%s

Focus on: financial risks, market risks, and operational concerns.`, data.Symbol, riskData)

	return ai.chat(prompt)
}

// GenerateOutlook provides future outlook
func (ai *AIClient) GenerateOutlook(data StockAnalysisData) (string, error) {
	prompt := fmt.Sprintf(`Based on this data for %s, provide a brief outlook (2-3 sentences):

%s

Be balanced and note this is not financial advice.`, data.Symbol, formatDataForAI(data))

	return ai.chat(prompt)
}

// SummarizeNews summarizes recent news sentiment
func (ai *AIClient) SummarizeNews(data *NewsSentimentResponse) (string, error) {
	if data == nil || len(data.Feed) == 0 {
		return "", fmt.Errorf("no news data")
	}

	newsData := formatNewsForAI(data)
	prompt := fmt.Sprintf(`Summarize the recent news sentiment (2-3 sentences):

%s

Focus on: overall sentiment, key themes, and potential market impact.`, newsData)

	return ai.chat(prompt)
}

// CustomAnalysis allows custom prompts with stock data
func (ai *AIClient) CustomAnalysis(data StockAnalysisData, customPrompt string) (string, error) {
	fullPrompt := fmt.Sprintf(`Stock: %s

Data:
%s

User Request: %s`, data.Symbol, formatDataForAI(data), customPrompt)

	return ai.chat(fullPrompt)
}


// Helper functions to format data for AI

func formatDataForAI(data StockAnalysisData) string {
	var sb strings.Builder

	// Price summary
	if data.Daily != nil && len(data.Daily.TimeSeries) > 0 {
		sb.WriteString(extractPriceSummary(data.Daily))
		sb.WriteString("\n\n")
	}

	// Earnings
	if data.Earnings != nil && len(data.Earnings.AnnualEarnings) > 0 {
		sb.WriteString("EARNINGS (Recent Years):\n")
		count := min(5, len(data.Earnings.AnnualEarnings))
		for i := 0; i < count; i++ {
			e := data.Earnings.AnnualEarnings[i]
			sb.WriteString(fmt.Sprintf("  %s: EPS $%s\n", e.FiscalDateEnding, e.ReportedEPS))
		}
		sb.WriteString("\n")
	}

	// Cash Flow
	if data.CashFlow != nil && len(data.CashFlow.AnnualReports) > 0 {
		r := data.CashFlow.AnnualReports[0]
		sb.WriteString(fmt.Sprintf("CASH FLOW (%s):\n", r.FiscalDateEnding))
		sb.WriteString(fmt.Sprintf("  Operating: %s\n", formatNum(r.OperatingCashflow)))
		sb.WriteString(fmt.Sprintf("  Investing: %s\n", formatNum(r.CashflowFromInvestment)))
		sb.WriteString(fmt.Sprintf("  Financing: %s\n", formatNum(r.CashflowFromFinancing)))
		sb.WriteString(fmt.Sprintf("  Net Income: %s\n\n", formatNum(r.NetIncome)))
	}

	// Balance Sheet
	if data.BalanceSheet != nil && len(data.BalanceSheet.AnnualReports) > 0 {
		r := data.BalanceSheet.AnnualReports[0]
		sb.WriteString(fmt.Sprintf("BALANCE SHEET (%s):\n", r.FiscalDateEnding))
		sb.WriteString(fmt.Sprintf("  Total Assets: %s\n", formatNum(r.TotalAssets)))
		sb.WriteString(fmt.Sprintf("  Total Liabilities: %s\n", formatNum(r.TotalLiabilities)))
		sb.WriteString(fmt.Sprintf("  Shareholder Equity: %s\n", formatNum(r.TotalShareholderEquity)))
		sb.WriteString(fmt.Sprintf("  Cash: %s\n", formatNum(r.CashAndCashEquivalentsAtCarryingValue)))
		sb.WriteString(fmt.Sprintf("  Long-term Debt: %s\n", formatNum(r.LongTermDebt)))
	}

	return sb.String()
}

func extractPriceSummary(data *TimeSeriesDailyResponse) string {
	if data == nil || len(data.TimeSeries) == 0 {
		return ""
	}

	// Sort dates
	var dates []string
	for d := range data.TimeSeries {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("PRICE DATA (%s):\n", data.MetaData.Symbol))

	// Latest price
	if len(dates) > 0 {
		latest := dates[len(dates)-1]
		latestData := data.TimeSeries[latest]
		sb.WriteString(fmt.Sprintf("  Latest (%s): Close $%s, Volume %s\n", latest, latestData.Close, latestData.Volume))
	}

	// Calculate stats
	var closes []float64
	var highs []float64
	var lows []float64
	for _, d := range dates {
		dp := data.TimeSeries[d]
		c, _ := strconv.ParseFloat(dp.Close, 64)
		h, _ := strconv.ParseFloat(dp.High, 64)
		l, _ := strconv.ParseFloat(dp.Low, 64)
		closes = append(closes, c)
		highs = append(highs, h)
		lows = append(lows, l)
	}

	if len(closes) > 0 {
		// Period high/low
		maxHigh := highs[0]
		minLow := lows[0]
		for _, h := range highs {
			if h > maxHigh {
				maxHigh = h
			}
		}
		for _, l := range lows {
			if l < minLow {
				minLow = l
			}
		}

		// Price change
		if len(closes) > 1 {
			first := closes[0]
			last := closes[len(closes)-1]
			change := ((last - first) / first) * 100
			sb.WriteString(fmt.Sprintf("  Period Change: %.2f%%\n", change))
		}

		sb.WriteString(fmt.Sprintf("  Period High: $%.2f\n", maxHigh))
		sb.WriteString(fmt.Sprintf("  Period Low: $%.2f\n", minLow))

		// Simple moving average
		if len(closes) >= 20 {
			sum := 0.0
			for i := len(closes) - 20; i < len(closes); i++ {
				sum += closes[i]
			}
			sma20 := sum / 20
			sb.WriteString(fmt.Sprintf("  20-day SMA: $%.2f\n", sma20))
		}
	}

	return sb.String()
}

func formatFundamentalsForAI(data StockAnalysisData) string {
	var sb strings.Builder

	// Earnings trend
	if data.Earnings != nil && len(data.Earnings.AnnualEarnings) >= 3 {
		sb.WriteString("EPS TREND:\n")
		for i := 0; i < min(5, len(data.Earnings.AnnualEarnings)); i++ {
			e := data.Earnings.AnnualEarnings[i]
			sb.WriteString(fmt.Sprintf("  %s: $%s\n", e.FiscalDateEnding, e.ReportedEPS))
		}
		sb.WriteString("\n")
	}

	// Cash flow health
	if data.CashFlow != nil && len(data.CashFlow.AnnualReports) > 0 {
		r := data.CashFlow.AnnualReports[0]
		sb.WriteString("CASH FLOW HEALTH:\n")
		sb.WriteString(fmt.Sprintf("  Operating CF: %s\n", formatNum(r.OperatingCashflow)))
		sb.WriteString(fmt.Sprintf("  CapEx: %s\n", formatNum(r.CapitalExpenditures)))
		sb.WriteString(fmt.Sprintf("  Dividends: %s\n\n", formatNum(r.DividendPayout)))
	}

	// Balance sheet ratios
	if data.BalanceSheet != nil && len(data.BalanceSheet.AnnualReports) > 0 {
		r := data.BalanceSheet.AnnualReports[0]
		assets, _ := strconv.ParseFloat(r.TotalAssets, 64)
		liabilities, _ := strconv.ParseFloat(r.TotalLiabilities, 64)
		equity, _ := strconv.ParseFloat(r.TotalShareholderEquity, 64)

		sb.WriteString("KEY RATIOS:\n")
		if equity > 0 {
			debtToEquity := liabilities / equity
			sb.WriteString(fmt.Sprintf("  Debt-to-Equity: %.2f\n", debtToEquity))
		}
		if assets > 0 {
			equityRatio := equity / assets
			sb.WriteString(fmt.Sprintf("  Equity Ratio: %.2f%%\n", equityRatio*100))
		}
	}

	return sb.String()
}

func formatRiskDataForAI(data StockAnalysisData) string {
	var sb strings.Builder

	// Debt levels
	if data.BalanceSheet != nil && len(data.BalanceSheet.AnnualReports) > 0 {
		r := data.BalanceSheet.AnnualReports[0]
		sb.WriteString(fmt.Sprintf("Debt: Long-term %s, Short-term %s\n", 
			formatNum(r.LongTermDebt), formatNum(r.ShortTermDebt)))
		sb.WriteString(fmt.Sprintf("Cash Position: %s\n", formatNum(r.CashAndCashEquivalentsAtCarryingValue)))
	}

	// Earnings volatility
	if data.Earnings != nil && len(data.Earnings.AnnualEarnings) >= 3 {
		var eps []float64
		for i := 0; i < min(5, len(data.Earnings.AnnualEarnings)); i++ {
			e, _ := strconv.ParseFloat(data.Earnings.AnnualEarnings[i].ReportedEPS, 64)
			eps = append(eps, e)
		}
		if len(eps) > 1 {
			// Check for declining trend
			declining := true
			for i := 1; i < len(eps); i++ {
				if eps[i] >= eps[i-1] {
					declining = false
					break
				}
			}
			if declining {
				sb.WriteString("Warning: Declining EPS trend\n")
			}
		}
	}

	// Price volatility
	if data.Daily != nil && len(data.Daily.TimeSeries) > 20 {
		sb.WriteString(fmt.Sprintf("Price data points: %d days\n", len(data.Daily.TimeSeries)))
	}

	return sb.String()
}

func formatNewsForAI(data *NewsSentimentResponse) string {
	if data == nil || len(data.Feed) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("RECENT NEWS:\n")

	count := min(5, len(data.Feed))
	for i := 0; i < count; i++ {
		item := data.Feed[i]
		sb.WriteString(fmt.Sprintf("- %s (Sentiment: %s, Score: %.2f)\n",
			truncate(item.Title, 80), item.OverallSentimentLabel, item.OverallSentimentScore))
	}

	return sb.String()
}

func formatNum(s string) string {
	if s == "" || s == "None" {
		return "N/A"
	}
	var num float64
	fmt.Sscanf(s, "%f", &num)

	negative := num < 0
	if negative {
		num = -num
	}

	var result string
	switch {
	case num >= 1e12:
		result = fmt.Sprintf("$%.2fT", num/1e12)
	case num >= 1e9:
		result = fmt.Sprintf("$%.2fB", num/1e9)
	case num >= 1e6:
		result = fmt.Sprintf("$%.2fM", num/1e6)
	default:
		result = fmt.Sprintf("$%.2f", num)
	}

	if negative {
		result = "-" + result
	}
	return result
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
