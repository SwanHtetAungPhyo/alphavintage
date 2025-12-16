package alphavintage

import "encoding/json"

// GetMarketStatus returns the current market status for major trading venues
func (c *Client) GetMarketStatus() (*MarketStatusResponse, error) {
	params := map[string]string{
		"function": "MARKET_STATUS",
	}

	body, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	var result MarketStatusResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
