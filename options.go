package zirael

import "time"

type Option func(*Client)

func WithHTTPTimeout(timeout time.Duration) Option{
	return func(client *Client) {
		client.timeout = timeout
	}
}

func WithHTTPClient(client IDoer) Option {
	return func(c *Client) {
		c.client = client
	}
}