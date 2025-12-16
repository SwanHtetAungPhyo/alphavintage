package alphavintage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const fdBaseURL = "https://api.financialdatasets.ai"

// FinancialDatasetsClient handles Financial Datasets API
type FinancialDatasetsClient struct {
	apiKey string
	resty  *resty.Client
}

// NewFinancialDatasetsClient creates a new Financial Datasets API client
func NewFinancialDatasetsClient(apiKey string) *FinancialDatasetsClient {
	return &FinancialDatasetsClient{
		apiKey: apiKey,
		resty:  resty.New().SetTimeout(30 * time.Second),
	}
}

func (c *FinancialDatasetsClient) doRequest(endpoint string, params map[string]string) ([]byte, error) {
	resp, err := c.resty.R().
		SetHeader("X-API-KEY", c.apiKey).
		SetQueryParams(params).
		Get(fdBaseURL + endpoint)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		var errResp struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		json.Unmarshal(resp.Body(), &errResp)
		if errResp.Message != "" {
			return nil, fmt.Errorf("API error: %s", errResp.Message)
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

// Period type for financial data
type FDPeriod string

const (
	FDPeriodAnnual    FDPeriod = "annual"
	FDPeriodQuarterly FDPeriod = "quarterly"
	FDPeriodTTM       FDPeriod = "ttm"
)

// Interval type for price data
type FDInterval string

const (
	FDIntervalSecond FDInterval = "second"
	FDIntervalMinute FDInterval = "minute"
	FDIntervalDay    FDInterval = "day"
	FDIntervalWeek   FDInterval = "week"
	FDIntervalMonth  FDInterval = "month"
	FDIntervalYear   FDInterval = "year"
)


// FD Types

// FDIncomeStatement represents income statement data
type FDIncomeStatement struct {
	Ticker                  string  `json:"ticker"`
	ReportPeriod            string  `json:"report_period"`
	FiscalPeriod            string  `json:"fiscal_period"`
	Period                  string  `json:"period"`
	Currency                string  `json:"currency"`
	Revenue                 float64 `json:"revenue"`
	CostOfRevenue           float64 `json:"cost_of_revenue"`
	GrossProfit             float64 `json:"gross_profit"`
	OperatingExpense        float64 `json:"operating_expense"`
	OperatingIncome         float64 `json:"operating_income"`
	InterestExpense         float64 `json:"interest_expense"`
	EBIT                    float64 `json:"ebit"`
	IncomeTaxExpense        float64 `json:"income_tax_expense"`
	NetIncome               float64 `json:"net_income"`
	EarningsPerShare        float64 `json:"earnings_per_share"`
	EarningsPerShareDiluted float64 `json:"earnings_per_share_diluted"`
	WeightedAverageShares   float64 `json:"weighted_average_shares"`
}

// FDBalanceSheet represents balance sheet data
type FDBalanceSheet struct {
	Ticker              string  `json:"ticker"`
	ReportPeriod        string  `json:"report_period"`
	FiscalPeriod        string  `json:"fiscal_period"`
	Period              string  `json:"period"`
	Currency            string  `json:"currency"`
	TotalAssets         float64 `json:"total_assets"`
	CurrentAssets       float64 `json:"current_assets"`
	CashAndEquivalents  float64 `json:"cash_and_equivalents"`
	Inventory           float64 `json:"inventory"`
	TotalLiabilities    float64 `json:"total_liabilities"`
	CurrentLiabilities  float64 `json:"current_liabilities"`
	CurrentDebt         float64 `json:"current_debt"`
	NonCurrentDebt      float64 `json:"non_current_debt"`
	TotalDebt           float64 `json:"total_debt"`
	ShareholdersEquity  float64 `json:"shareholders_equity"`
	RetainedEarnings    float64 `json:"retained_earnings"`
	OutstandingShares   float64 `json:"outstanding_shares"`
}

// FDCashFlowStatement represents cash flow data
type FDCashFlowStatement struct {
	Ticker                    string  `json:"ticker"`
	ReportPeriod              string  `json:"report_period"`
	FiscalPeriod              string  `json:"fiscal_period"`
	Period                    string  `json:"period"`
	Currency                  string  `json:"currency"`
	NetIncome                 float64 `json:"net_income"`
	DepreciationAmortization  float64 `json:"depreciation_and_amortization"`
	NetCashFlowFromOperations float64 `json:"net_cash_flow_from_operations"`
	CapitalExpenditure        float64 `json:"capital_expenditure"`
	NetCashFlowFromInvesting  float64 `json:"net_cash_flow_from_investing"`
	NetCashFlowFromFinancing  float64 `json:"net_cash_flow_from_financing"`
	FreeCashFlow              float64 `json:"free_cash_flow"`
	EndingCashBalance         float64 `json:"ending_cash_balance"`
}

// FDCompanyFacts represents company information
type FDCompanyFacts struct {
	Ticker            string  `json:"ticker"`
	Name              string  `json:"name"`
	CIK               string  `json:"cik"`
	Industry          string  `json:"industry"`
	Sector            string  `json:"sector"`
	Exchange          string  `json:"exchange"`
	IsActive          bool    `json:"is_active"`
	ListingDate       string  `json:"listing_date"`
	Location          string  `json:"location"`
	MarketCap         float64 `json:"market_cap"`
	NumberOfEmployees float64 `json:"number_of_employees"`
	WebsiteURL        string  `json:"website_url"`
}

// FDPrice represents price data
type FDPrice struct {
	Open             float64 `json:"open"`
	Close            float64 `json:"close"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	Volume           int64   `json:"volume"`
	Time             string  `json:"time"`
	TimeMilliseconds int64   `json:"time_milliseconds"`
}

// FDPriceSnapshot represents real-time price
type FDPriceSnapshot struct {
	Price            float64 `json:"price"`
	Ticker           string  `json:"ticker"`
	DayChange        float64 `json:"day_change"`
	DayChangePercent float64 `json:"day_change_percent"`
	MarketCap        float64 `json:"market_cap"`
	Time             string  `json:"time"`
}

// FDInsiderTrade represents insider trading data
type FDInsiderTrade struct {
	Ticker                       string  `json:"ticker"`
	Issuer                       string  `json:"issuer"`
	Name                         string  `json:"name"`
	Title                        string  `json:"title"`
	IsBoardDirector              bool    `json:"is_board_director"`
	TransactionDate              string  `json:"transaction_date"`
	TransactionShares            float64 `json:"transaction_shares"`
	TransactionPricePerShare     float64 `json:"transaction_price_per_share"`
	TransactionValue             float64 `json:"transaction_value"`
	SharesOwnedBeforeTransaction float64 `json:"shares_owned_before_transaction"`
	SharesOwnedAfterTransaction  float64 `json:"shares_owned_after_transaction"`
	FilingDate                   string  `json:"filing_date"`
}

// FDInstitutionalOwnership represents institutional holdings
type FDInstitutionalOwnership struct {
	Ticker       string  `json:"ticker"`
	Investor     string  `json:"investor"`
	ReportPeriod string  `json:"report_period"`
	Price        float64 `json:"price"`
	Shares       float64 `json:"shares"`
	MarketValue  float64 `json:"market_value"`
}

// FDNews represents news article
type FDNews struct {
	Ticker    string `json:"ticker"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Source    string `json:"source"`
	Date      string `json:"date"`
	URL       string `json:"url"`
	Sentiment string `json:"sentiment"`
}

// FDFinancialMetrics represents financial ratios
type FDFinancialMetrics struct {
	Ticker                      string  `json:"ticker"`
	MarketCap                   float64 `json:"market_cap"`
	EnterpriseValue             float64 `json:"enterprise_value"`
	PriceToEarningsRatio        float64 `json:"price_to_earnings_ratio"`
	PriceToBookRatio            float64 `json:"price_to_book_ratio"`
	PriceToSalesRatio           float64 `json:"price_to_sales_ratio"`
	EVToEBITDA                  float64 `json:"enterprise_value_to_ebitda_ratio"`
	GrossMargin                 float64 `json:"gross_margin"`
	OperatingMargin             float64 `json:"operating_margin"`
	NetMargin                   float64 `json:"net_margin"`
	ReturnOnEquity              float64 `json:"return_on_equity"`
	ReturnOnAssets              float64 `json:"return_on_assets"`
	CurrentRatio                float64 `json:"current_ratio"`
	QuickRatio                  float64 `json:"quick_ratio"`
	DebtToEquity                float64 `json:"debt_to_equity"`
	DebtToAssets                float64 `json:"debt_to_assets"`
	RevenueGrowth               float64 `json:"revenue_growth"`
	EarningsGrowth              float64 `json:"earnings_growth"`
	EarningsPerShare            float64 `json:"earnings_per_share"`
	BookValuePerShare           float64 `json:"book_value_per_share"`
	FreeCashFlowPerShare        float64 `json:"free_cash_flow_per_share"`
}


// API Methods

// GetIncomeStatements returns income statements for a ticker
func (c *FinancialDatasetsClient) GetIncomeStatements(ticker string, period FDPeriod, limit int) ([]FDIncomeStatement, error) {
	params := map[string]string{
		"ticker": ticker,
		"period": string(period),
	}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}

	body, err := c.doRequest("/financials/income-statements", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		IncomeStatements []FDIncomeStatement `json:"income_statements"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.IncomeStatements, nil
}

// GetBalanceSheets returns balance sheets for a ticker
func (c *FinancialDatasetsClient) GetBalanceSheets(ticker string, period FDPeriod, limit int) ([]FDBalanceSheet, error) {
	params := map[string]string{
		"ticker": ticker,
		"period": string(period),
	}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}

	body, err := c.doRequest("/financials/balance-sheets", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		BalanceSheets []FDBalanceSheet `json:"balance_sheets"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.BalanceSheets, nil
}

// GetCashFlowStatements returns cash flow statements for a ticker
func (c *FinancialDatasetsClient) GetCashFlowStatements(ticker string, period FDPeriod, limit int) ([]FDCashFlowStatement, error) {
	params := map[string]string{
		"ticker": ticker,
		"period": string(period),
	}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}

	body, err := c.doRequest("/financials/cash-flow-statements", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		CashFlowStatements []FDCashFlowStatement `json:"cash_flow_statements"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.CashFlowStatements, nil
}

// GetCompanyFacts returns company information
func (c *FinancialDatasetsClient) GetCompanyFacts(ticker string) (*FDCompanyFacts, error) {
	params := map[string]string{"ticker": ticker}

	body, err := c.doRequest("/company/facts", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		CompanyFacts FDCompanyFacts `json:"company_facts"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp.CompanyFacts, nil
}

// GetPriceSnapshot returns real-time price
func (c *FinancialDatasetsClient) GetPriceSnapshot(ticker string) (*FDPriceSnapshot, error) {
	params := map[string]string{"ticker": ticker}

	body, err := c.doRequest("/prices/snapshot", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Snapshot FDPriceSnapshot `json:"snapshot"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp.Snapshot, nil
}

// GetPrices returns historical price data
func (c *FinancialDatasetsClient) GetPrices(ticker string, interval FDInterval, multiplier int, startDate, endDate string, limit int) ([]FDPrice, error) {
	params := map[string]string{
		"ticker":              ticker,
		"interval":            string(interval),
		"interval_multiplier": fmt.Sprintf("%d", multiplier),
		"start_date":          startDate,
		"end_date":            endDate,
	}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}

	body, err := c.doRequest("/prices", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Prices []FDPrice `json:"prices"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Prices, nil
}

// GetInsiderTrades returns insider trading data
func (c *FinancialDatasetsClient) GetInsiderTrades(ticker string, limit int) ([]FDInsiderTrade, error) {
	params := map[string]string{"ticker": ticker}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}

	body, err := c.doRequest("/insider-trades", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		InsiderTrades []FDInsiderTrade `json:"insider_trades"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.InsiderTrades, nil
}

// GetInstitutionalOwnership returns institutional holdings
func (c *FinancialDatasetsClient) GetInstitutionalOwnership(ticker string, limit int) ([]FDInstitutionalOwnership, error) {
	params := map[string]string{"ticker": ticker}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}

	body, err := c.doRequest("/institutional-ownership", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Ownership []FDInstitutionalOwnership `json:"institutional-ownership"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Ownership, nil
}

// GetNews returns news articles
func (c *FinancialDatasetsClient) GetNews(ticker string, startDate, endDate string, limit int) ([]FDNews, error) {
	params := map[string]string{"ticker": ticker}
	if startDate != "" {
		params["start_date"] = startDate
	}
	if endDate != "" {
		params["end_date"] = endDate
	}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}

	body, err := c.doRequest("/news", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		News []FDNews `json:"news"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.News, nil
}

// GetFinancialMetrics returns financial ratios and metrics
func (c *FinancialDatasetsClient) GetFinancialMetrics(ticker string, period FDPeriod, limit int) ([]FDFinancialMetrics, error) {
	params := map[string]string{
		"ticker": ticker,
		"period": string(period),
	}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}

	body, err := c.doRequest("/financial-metrics", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Metrics []FDFinancialMetrics `json:"financial_metrics"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Metrics, nil
}

// GetFinancialMetricsSnapshot returns current financial metrics
func (c *FinancialDatasetsClient) GetFinancialMetricsSnapshot(ticker string) (*FDFinancialMetrics, error) {
	params := map[string]string{"ticker": ticker}

	body, err := c.doRequest("/financial-metrics/snapshot", params)
	if err != nil {
		return nil, err
	}

	var resp FDFinancialMetrics
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
