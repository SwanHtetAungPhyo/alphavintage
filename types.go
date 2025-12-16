package alphavintage

// MarketStatusResponse represents the market status API response
type MarketStatusResponse struct {
	Endpoint string   `json:"endpoint"`
	Markets  []Market `json:"markets"`
}

// Market represents a single market's status
type Market struct {
	MarketType       string `json:"market_type"`
	Region           string `json:"region"`
	PrimaryExchanges string `json:"primary_exchanges"`
	LocalOpen        string `json:"local_open"`
	LocalClose       string `json:"local_close"`
	CurrentStatus    string `json:"current_status"`
	Notes            string `json:"notes"`
}

// TimeSeriesDailyResponse represents daily time series data
type TimeSeriesDailyResponse struct {
	MetaData   TimeSeriesMetaData        `json:"Meta Data"`
	TimeSeries map[string]DailyDataPoint `json:"Time Series (Daily)"`
}

// TimeSeriesMetaData contains metadata for time series
type TimeSeriesMetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
}

// DailyDataPoint represents OHLCV data for a single day
type DailyDataPoint struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

// TimeSeriesIntradayResponse represents intraday time series data
type TimeSeriesIntradayResponse struct {
	MetaData   IntradayMetaData             `json:"Meta Data"`
	TimeSeries map[string]IntradayDataPoint `json:"Time Series (5min)"`
}

// IntradayMetaData contains metadata for intraday series
type IntradayMetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	Interval      string `json:"4. Interval"`
	OutputSize    string `json:"5. Output Size"`
	TimeZone      string `json:"6. Time Zone"`
}

// IntradayDataPoint represents OHLCV data for intraday
type IntradayDataPoint struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}


// NewsSentimentResponse represents news sentiment API response
type NewsSentimentResponse struct {
	Items                    string        `json:"items"`
	SentimentScoreDefinition string        `json:"sentiment_score_definition"`
	RelevanceScoreDefinition string        `json:"relevance_score_definition"`
	Feed                     []NewsFeedItem `json:"feed"`
}

// NewsFeedItem represents a single news item
type NewsFeedItem struct {
	Title                 string            `json:"title"`
	URL                   string            `json:"url"`
	TimePublished         string            `json:"time_published"`
	Authors               []string          `json:"authors"`
	Summary               string            `json:"summary"`
	BannerImage           string            `json:"banner_image"`
	Source                string            `json:"source"`
	CategoryWithinSource  string            `json:"category_within_source"`
	SourceDomain          string            `json:"source_domain"`
	Topics                []Topic           `json:"topics"`
	OverallSentimentScore float64           `json:"overall_sentiment_score"`
	OverallSentimentLabel string            `json:"overall_sentiment_label"`
	TickerSentiment       []TickerSentiment `json:"ticker_sentiment"`
}

// Topic represents a news topic
type Topic struct {
	Topic          string `json:"topic"`
	RelevanceScore string `json:"relevance_score"`
}

// TickerSentiment represents sentiment for a specific ticker
type TickerSentiment struct {
	Ticker               string `json:"ticker"`
	RelevanceScore       string `json:"relevance_score"`
	TickerSentimentScore string `json:"ticker_sentiment_score"`
	TickerSentimentLabel string `json:"ticker_sentiment_label"`
}

// BalanceSheetResponse represents balance sheet API response
type BalanceSheetResponse struct {
	Symbol           string               `json:"symbol"`
	AnnualReports    []BalanceSheetReport `json:"annualReports"`
	QuarterlyReports []BalanceSheetReport `json:"quarterlyReports"`
}

// BalanceSheetReport represents a single balance sheet report
type BalanceSheetReport struct {
	FiscalDateEnding                       string `json:"fiscalDateEnding"`
	ReportedCurrency                       string `json:"reportedCurrency"`
	TotalAssets                            string `json:"totalAssets"`
	TotalCurrentAssets                     string `json:"totalCurrentAssets"`
	CashAndCashEquivalentsAtCarryingValue  string `json:"cashAndCashEquivalentsAtCarryingValue"`
	CashAndShortTermInvestments            string `json:"cashAndShortTermInvestments"`
	Inventory                              string `json:"inventory"`
	CurrentNetReceivables                  string `json:"currentNetReceivables"`
	TotalNonCurrentAssets                  string `json:"totalNonCurrentAssets"`
	PropertyPlantEquipment                 string `json:"propertyPlantEquipment"`
	AccumulatedDepreciationAmortizationPPE string `json:"accumulatedDepreciationAmortizationPPE"`
	IntangibleAssets                       string `json:"intangibleAssets"`
	IntangibleAssetsExcludingGoodwill      string `json:"intangibleAssetsExcludingGoodwill"`
	Goodwill                               string `json:"goodwill"`
	Investments                            string `json:"investments"`
	LongTermInvestments                    string `json:"longTermInvestments"`
	ShortTermInvestments                   string `json:"shortTermInvestments"`
	OtherCurrentAssets                     string `json:"otherCurrentAssets"`
	OtherNonCurrentAssets                  string `json:"otherNonCurrentAssets"`
	TotalLiabilities                       string `json:"totalLiabilities"`
	TotalCurrentLiabilities                string `json:"totalCurrentLiabilities"`
	CurrentAccountsPayable                 string `json:"currentAccountsPayable"`
	DeferredRevenue                        string `json:"deferredRevenue"`
	CurrentDebt                            string `json:"currentDebt"`
	ShortTermDebt                          string `json:"shortTermDebt"`
	TotalNonCurrentLiabilities             string `json:"totalNonCurrentLiabilities"`
	CapitalLeaseObligations                string `json:"capitalLeaseObligations"`
	LongTermDebt                           string `json:"longTermDebt"`
	CurrentLongTermDebt                    string `json:"currentLongTermDebt"`
	LongTermDebtNoncurrent                 string `json:"longTermDebtNoncurrent"`
	ShortLongTermDebtTotal                 string `json:"shortLongTermDebtTotal"`
	OtherCurrentLiabilities                string `json:"otherCurrentLiabilities"`
	OtherNonCurrentLiabilities             string `json:"otherNonCurrentLiabilities"`
	TotalShareholderEquity                 string `json:"totalShareholderEquity"`
	TreasuryStock                          string `json:"treasuryStock"`
	RetainedEarnings                       string `json:"retainedEarnings"`
	CommonStock                            string `json:"commonStock"`
	CommonStockSharesOutstanding           string `json:"commonStockSharesOutstanding"`
}


// CashFlowResponse represents cash flow API response
type CashFlowResponse struct {
	Symbol           string           `json:"symbol"`
	AnnualReports    []CashFlowReport `json:"annualReports"`
	QuarterlyReports []CashFlowReport `json:"quarterlyReports"`
}

// CashFlowReport represents a single cash flow report
type CashFlowReport struct {
	FiscalDateEnding                                          string `json:"fiscalDateEnding"`
	ReportedCurrency                                          string `json:"reportedCurrency"`
	OperatingCashflow                                         string `json:"operatingCashflow"`
	PaymentsForOperatingActivities                            string `json:"paymentsForOperatingActivities"`
	ProceedsFromOperatingActivities                           string `json:"proceedsFromOperatingActivities"`
	ChangeInOperatingLiabilities                              string `json:"changeInOperatingLiabilities"`
	ChangeInOperatingAssets                                   string `json:"changeInOperatingAssets"`
	DepreciationDepletionAndAmortization                      string `json:"depreciationDepletionAndAmortization"`
	CapitalExpenditures                                       string `json:"capitalExpenditures"`
	ChangeInReceivables                                       string `json:"changeInReceivables"`
	ChangeInInventory                                         string `json:"changeInInventory"`
	ProfitLoss                                                string `json:"profitLoss"`
	CashflowFromInvestment                                    string `json:"cashflowFromInvestment"`
	CashflowFromFinancing                                     string `json:"cashflowFromFinancing"`
	ProceedsFromRepaymentsOfShortTermDebt                     string `json:"proceedsFromRepaymentsOfShortTermDebt"`
	PaymentsForRepurchaseOfCommonStock                        string `json:"paymentsForRepurchaseOfCommonStock"`
	PaymentsForRepurchaseOfEquity                             string `json:"paymentsForRepurchaseOfEquity"`
	PaymentsForRepurchaseOfPreferredStock                     string `json:"paymentsForRepurchaseOfPreferredStock"`
	DividendPayout                                            string `json:"dividendPayout"`
	DividendPayoutCommonStock                                 string `json:"dividendPayoutCommonStock"`
	DividendPayoutPreferredStock                              string `json:"dividendPayoutPreferredStock"`
	ProceedsFromIssuanceOfCommonStock                         string `json:"proceedsFromIssuanceOfCommonStock"`
	ProceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet string `json:"proceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet"`
	ProceedsFromIssuanceOfPreferredStock                      string `json:"proceedsFromIssuanceOfPreferredStock"`
	ProceedsFromRepurchaseOfEquity                            string `json:"proceedsFromRepurchaseOfEquity"`
	ProceedsFromSaleOfTreasuryStock                           string `json:"proceedsFromSaleOfTreasuryStock"`
	ChangeInCashAndCashEquivalents                            string `json:"changeInCashAndCashEquivalents"`
	ChangeInExchangeRate                                      string `json:"changeInExchangeRate"`
	NetIncome                                                 string `json:"netIncome"`
}

// EarningsResponse represents earnings API response
type EarningsResponse struct {
	Symbol             string            `json:"symbol"`
	AnnualEarnings     []AnnualEarning   `json:"annualEarnings"`
	QuarterlyEarnings  []QuarterlyEarning `json:"quarterlyEarnings"`
}

// AnnualEarning represents annual earnings data
type AnnualEarning struct {
	FiscalDateEnding string `json:"fiscalDateEnding"`
	ReportedEPS      string `json:"reportedEPS"`
}

// QuarterlyEarning represents quarterly earnings data
type QuarterlyEarning struct {
	FiscalDateEnding   string `json:"fiscalDateEnding"`
	ReportedDate       string `json:"reportedDate"`
	ReportedEPS        string `json:"reportedEPS"`
	EstimatedEPS       string `json:"estimatedEPS"`
	Surprise           string `json:"surprise"`
	SurprisePercentage string `json:"surprisePercentage"`
}
