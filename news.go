package alphavintage

import (
	"encoding/json"
	"strconv"
)

// NewsSentimentOptions contains options for news sentiment API
type NewsSentimentOptions struct {
	Tickers   string // Comma-separated ticker symbols
	Topics    string // Comma-separated topics
	TimeFrom  string // YYYYMMDDTHHMM format
	TimeTo    string // YYYYMMDDTHHMM format
	Sort      string // LATEST, EARLIEST, RELEVANCE
	Limit     int    // Number of results (max 1000)
}

// GetNewsSentiment returns news and sentiment data
func (c *Client) GetNewsSentiment(opts *NewsSentimentOptions) (*NewsSentimentResponse, error) {
	params := map[string]string{
		"function": "NEWS_SENTIMENT",
	}

	if opts != nil {
		if opts.Tickers != "" {
			params["tickers"] = opts.Tickers
		}
		if opts.Topics != "" {
			params["topics"] = opts.Topics
		}
		if opts.TimeFrom != "" {
			params["time_from"] = opts.TimeFrom
		}
		if opts.TimeTo != "" {
			params["time_to"] = opts.TimeTo
		}
		if opts.Sort != "" {
			params["sort"] = opts.Sort
		}
		if opts.Limit > 0 {
			params["limit"] = strconv.Itoa(opts.Limit)
		}
	}

	body, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	var result NewsSentimentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
