package alphavintage

import (
	"encoding/json"
	"fmt"
)

// OutputSize represents the output size option
type OutputSize string

const (
	OutputSizeCompact OutputSize = "compact"
	OutputSizeFull    OutputSize = "full"
)

// Interval represents intraday interval options
type Interval string

const (
	Interval1Min  Interval = "1min"
	Interval5Min  Interval = "5min"
	Interval15Min Interval = "15min"
	Interval30Min Interval = "30min"
	Interval60Min Interval = "60min"
)

// GetTimeSeriesDaily returns daily OHLCV data for a symbol
func (c *Client) GetTimeSeriesDaily(symbol string, outputSize OutputSize) (*TimeSeriesDailyResponse, error) {
	params := map[string]string{
		"function": "TIME_SERIES_DAILY",
		"symbol":   symbol,
	}
	if outputSize != "" {
		params["outputsize"] = string(outputSize)
	}

	body, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	var result TimeSeriesDailyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTimeSeriesIntraday returns intraday OHLCV data for a symbol
func (c *Client) GetTimeSeriesIntraday(symbol string, interval Interval, outputSize OutputSize) (*TimeSeriesIntradayResponse, error) {
	params := map[string]string{
		"function": "TIME_SERIES_INTRADAY",
		"symbol":   symbol,
		"interval": string(interval),
	}
	if outputSize != "" {
		params["outputsize"] = string(outputSize)
	}

	body, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	// Dynamic key based on interval
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	var result TimeSeriesIntradayResponse
	if err := json.Unmarshal(raw["Meta Data"], &result.MetaData); err != nil {
		return nil, err
	}

	timeSeriesKey := fmt.Sprintf("Time Series (%s)", interval)
	if tsData, ok := raw[timeSeriesKey]; ok {
		result.TimeSeries = make(map[string]IntradayDataPoint)
		if err := json.Unmarshal(tsData, &result.TimeSeries); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

// FilterIntradayByDate filters intraday data for a specific date (YYYY-MM-DD)
func FilterIntradayByDate(data *TimeSeriesIntradayResponse, date string) *TimeSeriesIntradayResponse {
	if data == nil {
		return nil
	}

	filtered := &TimeSeriesIntradayResponse{
		MetaData:   data.MetaData,
		TimeSeries: make(map[string]IntradayDataPoint),
	}

	for timestamp, point := range data.TimeSeries {
		// Timestamp format: "2024-12-16 10:30:00"
		if len(timestamp) >= 10 && timestamp[:10] == date {
			filtered.TimeSeries[timestamp] = point
		}
	}

	return filtered
}

// GetSingleDayData returns intraday data for a specific date
// Note: Alpha Vantage free tier only returns recent data (last 1-2 trading days)
// Premium subscription required for extended intraday history
func (c *Client) GetSingleDayData(symbol string, date string, interval Interval) (*TimeSeriesIntradayResponse, error) {
	// Fetch intraday data
	data, err := c.GetTimeSeriesIntraday(symbol, interval, OutputSizeFull)
	if err != nil {
		return nil, err
	}

	// Filter for the specific date
	filtered := FilterIntradayByDate(data, date)

	if len(filtered.TimeSeries) == 0 {
		return nil, fmt.Errorf("no data available for %s on %s (Alpha Vantage free tier only provides recent 1-2 days)", symbol, date)
	}

	return filtered, nil
}

// GetDailyDataForDate returns a single day's OHLCV from daily time series
func (c *Client) GetDailyDataForDate(symbol string, date string) (*DailyDataPoint, error) {
	data, err := c.GetTimeSeriesDaily(symbol, OutputSizeCompact)
	if err != nil {
		return nil, err
	}

	if point, ok := data.TimeSeries[date]; ok {
		return &point, nil
	}

	return nil, fmt.Errorf("no data for %s on %s", symbol, date)
}

// GetIntradaySummary returns summary statistics for intraday data
type IntradaySummary struct {
	Symbol     string
	Date       string
	Open       float64 // First price of the day
	High       float64 // Highest price
	Low        float64 // Lowest price
	Close      float64 // Last price of the day
	TotalVol   int64   // Total volume
	DataPoints int     // Number of data points
	Interval   string
}

// GetIntradaySummary calculates summary for intraday data
func GetIntradaySummary(data *TimeSeriesIntradayResponse) (*IntradaySummary, error) {
	if data == nil || len(data.TimeSeries) == 0 {
		return nil, fmt.Errorf("no data")
	}

	summary := &IntradaySummary{
		Symbol:     data.MetaData.Symbol,
		Interval:   data.MetaData.Interval,
		DataPoints: len(data.TimeSeries),
	}

	// Sort timestamps to get proper open/close
	type timePoint struct {
		time  string
		point IntradayDataPoint
	}
	var points []timePoint
	for t, p := range data.TimeSeries {
		points = append(points, timePoint{t, p})
	}

	// Sort by time
	for i := 0; i < len(points)-1; i++ {
		for j := i + 1; j < len(points); j++ {
			if points[i].time > points[j].time {
				points[i], points[j] = points[j], points[i]
			}
		}
	}

	if len(points) > 0 {
		summary.Date = points[0].time[:10]
	}

	first := true
	for _, tp := range points {
		p := tp.point
		open, _ := parseFloat(p.Open)
		high, _ := parseFloat(p.High)
		low, _ := parseFloat(p.Low)
		close, _ := parseFloat(p.Close)
		vol, _ := parseInt(p.Volume)

		if first {
			summary.Open = open
			summary.High = high
			summary.Low = low
			first = false
		}

		if high > summary.High {
			summary.High = high
		}
		if low < summary.Low {
			summary.Low = low
		}
		summary.Close = close
		summary.TotalVol += vol
	}

	return summary, nil
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

func parseInt(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// FilterDailyByDateRange filters daily data for a date range (inclusive)
// startDate and endDate format: "YYYY-MM-DD"
// Pass empty string for startDate to get all data up to endDate
// Pass empty string for endDate to get all data from startDate onwards
func FilterDailyByDateRange(data *TimeSeriesDailyResponse, startDate, endDate string) *TimeSeriesDailyResponse {
	if data == nil {
		return nil
	}

	filtered := &TimeSeriesDailyResponse{
		MetaData:   data.MetaData,
		TimeSeries: make(map[string]DailyDataPoint),
	}

	for date, point := range data.TimeSeries {
		// Check if date is within range
		if startDate != "" && date < startDate {
			continue
		}
		if endDate != "" && date > endDate {
			continue
		}
		filtered.TimeSeries[date] = point
	}

	return filtered
}

// FilterDailyLastNDays filters daily data for the last N trading days
func FilterDailyLastNDays(data *TimeSeriesDailyResponse, days int) *TimeSeriesDailyResponse {
	if data == nil || days <= 0 {
		return nil
	}

	// Get sorted dates
	dates := GetSortedDates(data)
	if len(dates) == 0 {
		return nil
	}

	// Take last N days
	if days > len(dates) {
		days = len(dates)
	}
	recentDates := dates[len(dates)-days:]

	filtered := &TimeSeriesDailyResponse{
		MetaData:   data.MetaData,
		TimeSeries: make(map[string]DailyDataPoint),
	}

	for _, date := range recentDates {
		if point, ok := data.TimeSeries[date]; ok {
			filtered.TimeSeries[date] = point
		}
	}

	return filtered
}

// GetSortedDates returns all dates from daily data sorted ascending
func GetSortedDates(data *TimeSeriesDailyResponse) []string {
	if data == nil {
		return nil
	}

	dates := make([]string, 0, len(data.TimeSeries))
	for date := range data.TimeSeries {
		dates = append(dates, date)
	}

	// Sort ascending
	for i := 0; i < len(dates)-1; i++ {
		for j := i + 1; j < len(dates); j++ {
			if dates[i] > dates[j] {
				dates[i], dates[j] = dates[j], dates[i]
			}
		}
	}

	return dates
}

// DailyRangeSummary contains summary statistics for a date range
type DailyRangeSummary struct {
	Symbol       string
	StartDate    string
	EndDate      string
	TradingDays  int
	PeriodOpen   float64 // First day's open
	PeriodHigh   float64 // Highest high in period
	PeriodLow    float64 // Lowest low in period
	PeriodClose  float64 // Last day's close
	TotalVolume  int64
	AvgVolume    int64
	PriceChange  float64 // Close - Open
	PriceChangePct float64
	HighDate     string // Date of highest price
	LowDate      string // Date of lowest price
}

// GetDailyRangeSummary calculates summary statistics for daily data
func GetDailyRangeSummary(data *TimeSeriesDailyResponse) (*DailyRangeSummary, error) {
	if data == nil || len(data.TimeSeries) == 0 {
		return nil, fmt.Errorf("no data")
	}

	dates := GetSortedDates(data)
	if len(dates) == 0 {
		return nil, fmt.Errorf("no dates")
	}

	summary := &DailyRangeSummary{
		Symbol:      data.MetaData.Symbol,
		StartDate:   dates[0],
		EndDate:     dates[len(dates)-1],
		TradingDays: len(dates),
	}

	first := true
	for _, date := range dates {
		point := data.TimeSeries[date]
		open, _ := parseFloat(point.Open)
		high, _ := parseFloat(point.High)
		low, _ := parseFloat(point.Low)
		close, _ := parseFloat(point.Close)
		vol, _ := parseInt(point.Volume)

		if first {
			summary.PeriodOpen = open
			summary.PeriodHigh = high
			summary.PeriodLow = low
			summary.HighDate = date
			summary.LowDate = date
			first = false
		}

		if high > summary.PeriodHigh {
			summary.PeriodHigh = high
			summary.HighDate = date
		}
		if low < summary.PeriodLow {
			summary.PeriodLow = low
			summary.LowDate = date
		}

		summary.PeriodClose = close
		summary.TotalVolume += vol
	}

	if summary.TradingDays > 0 {
		summary.AvgVolume = summary.TotalVolume / int64(summary.TradingDays)
	}

	summary.PriceChange = summary.PeriodClose - summary.PeriodOpen
	if summary.PeriodOpen != 0 {
		summary.PriceChangePct = (summary.PriceChange / summary.PeriodOpen) * 100
	}

	return summary, nil
}

// GetDailyDataPoint returns a single day's data from already-fetched daily response
func GetDailyDataPoint(data *TimeSeriesDailyResponse, date string) (*DailyDataPoint, bool) {
	if data == nil {
		return nil, false
	}
	point, ok := data.TimeSeries[date]
	return &point, ok
}

// GetMostRecentDate returns the most recent trading date in the data
func GetMostRecentDate(data *TimeSeriesDailyResponse) string {
	if data == nil || len(data.TimeSeries) == 0 {
		return ""
	}
	dates := GetSortedDates(data)
	return dates[len(dates)-1]
}

// GetOldestDate returns the oldest trading date in the data
func GetOldestDate(data *TimeSeriesDailyResponse) string {
	if data == nil || len(data.TimeSeries) == 0 {
		return ""
	}
	dates := GetSortedDates(data)
	return dates[0]
}
