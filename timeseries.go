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
