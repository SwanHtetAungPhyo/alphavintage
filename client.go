package alphavintage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const baseURL = "https://www.alphavantage.co/query"

// Client is the Alpha Vantage API client
type Client struct {
	apiKey string
	resty  *resty.Client
}

// NewClient creates a new Alpha Vantage client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		resty:  resty.New().SetTimeout(30 * time.Second),
	}
}

// WithRestyClient sets a custom resty client
func (c *Client) WithRestyClient(client *resty.Client) *Client {
	c.resty = client
	return c
}

func (c *Client) doRequest(params map[string]string) ([]byte, error) {
	params["apikey"] = c.apiKey

	resp, err := c.resty.R().SetQueryParams(params).Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	body := resp.Body()

	// Check for API error response
	var apiErr struct {
		ErrorMessage string `json:"Error Message"`
		Note         string `json:"Note"`
		Information  string `json:"Information"`
	}
	if json.Unmarshal(body, &apiErr) == nil {
		if apiErr.ErrorMessage != "" {
			return nil, fmt.Errorf("API error: %s", apiErr.ErrorMessage)
		}
		if apiErr.Note != "" {
			return nil, fmt.Errorf("API rate limit: %s", apiErr.Note)
		}
		if apiErr.Information != "" {
			return nil, fmt.Errorf("API info: %s", apiErr.Information)
		}
	}

	return body, nil
}
