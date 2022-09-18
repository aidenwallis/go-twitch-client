package client

import (
	"net/http"
	"time"
)

// Client is a wrapper around the http client
type Client struct {
	*http.Client
}

// Options are the base client-level options
type Options struct {
	// RequestTimeout defines the client-level request timeout. Leave empty for no
	// timeout.
	RequestTimeout time.Duration

	// Transport defines a HTTP transport, useful for defining keep-alive settings,
	// timeouts, or mocking client behaviour in tests. In high throughput or advanced
	// scenarios, it's likely that you will want to define custom transport settings to
	// properly adjust your connection pooling and timeout logic accordingly.
	Transport http.RoundTripper
}

// NewClient creates a new instance of client.
func NewClient(options *Options) *Client {
	tr := http.DefaultTransport
	if options.Transport != nil {
		tr = options.Transport
	}

	return &Client{
		Client: &http.Client{
			Timeout:   options.RequestTimeout,
			Transport: tr,
		},
	}
}
