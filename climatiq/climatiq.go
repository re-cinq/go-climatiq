package climatiq

import (
	"net/http"
	"net/url"
)

const (
	defaultBaseURL   = "https://beta4.api.climatiq.io/"
	defaultUserAgent = "go-climatiq"
)

type clientOpts func(*Client)

// Client is the structure use to communicate with the climatiq API
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// base URL for API requests.
	// defaults to the public climatiq api
	baseURL *url.URL

	// user agent used to make requests
	userAgent string

	//token is used for authentication to the climatiq API
	token string
}

// NewClient returns an instantiated instance of a client
// with the ability to override values with various options
func NewClient(opts ...clientOpts) *Client {
	url, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:  &http.Client{},
		baseURL: url,
	}

	return c
}

// WithBaseURL is an option to overwrite the BaseURL
func WithBaseURL(u string) clientOpts {
	url, _ := url.Parse(u)
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithUserAgent is an option to overwrite the UserAgent
func WithUserAgent(u string) clientOpts {
	return func(c *Client) {
		c.userAgent = u
	}
}

// WithClient is an option to overwrite the default client
func WithClient(cli *http.Client) clientOpts {
	return func(c *Client) {
		c.client = cli
	}
}

// Do is used to make the actual http requests
func (c *Client) Do(r *http.Request) {}
