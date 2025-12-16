package alphavintage

import "encoding/json"

// GetBalanceSheet returns balance sheet data for a symbol
func (c *Client) GetBalanceSheet(symbol string) (*BalanceSheetResponse, error) {
	params := map[string]string{
		"function": "BALANCE_SHEET",
		"symbol":   symbol,
	}

	body, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	var result BalanceSheetResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCashFlow returns cash flow data for a symbol
func (c *Client) GetCashFlow(symbol string) (*CashFlowResponse, error) {
	params := map[string]string{
		"function": "CASH_FLOW",
		"symbol":   symbol,
	}

	body, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	var result CashFlowResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetEarnings returns earnings data for a symbol
func (c *Client) GetEarnings(symbol string) (*EarningsResponse, error) {
	params := map[string]string{
		"function": "EARNINGS",
		"symbol":   symbol,
	}

	body, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	var result EarningsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
