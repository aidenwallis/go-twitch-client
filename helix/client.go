package helix

import (
	"context"
	"net/http"
	"time"

	"github.com/aidenwallis/go-twitch-client/internal/client"
)

// Client defines the helix client
type Client interface {
	Ads
	Users
}

// Pagination is the helix pagination response
type Pagination struct {
	// Cursor value, to be used in a subsequent request to specify the starting point of the next set of results.
	Cursor string `json:"cursor"`
}

type helixClient struct {
	*client.Client

	accessTokenLoader AccessTokenLoader
	clientID          string
}

// RequestOptions are the common options passed to every request
type RequestOptions struct {
	// Token is the OAuth bearer token for the Twitch request
	Token string
}

// AccessTokenLoader is the loader to return a generic app access token
type AccessTokenLoader func(ctx context.Context) (appAccessToken string, err error)

// ClientOptions defines all options this client supports.
type ClientOptions struct {
	// ClientID is your third-party client ID that you will use for interacting with the
	// Twitch API. You can obtain one at dev.twitch.tv.
	//
	// Your bearer tokens must have been created using this client ID, else Helix will
	// return an error.
	ClientID string

	// RequestTimeout defines a client-level request timeout duration for your client.
	// You may also cancel individual requests by using context cancellation.
	//
	// It's strongly recommended that you define a request timeout, but you can leave
	// empty to provide no request timeout.
	RequestTimeout time.Duration

	// Transport defines a HTTP transport, useful for defining keep-alive settings,
	// timeouts, or mocking client behaviour in tests. In high throughput or advanced
	// scenarios, it's likely that you will want to define custom transport settings to
	// properly adjust your connection pooling and timeout logic accordingly.
	Transport http.RoundTripper

	// AccessTokenLoader is an optional argument that can be invoked if you would like
	// requests that do not provide a bearer token to have a fallback mechanism.
	//
	// For example, you could use this for fetching an app access token, and falling back
	// to that in these cases.
	//
	// If you do not define a loader, it will simply return an error in those cases instead.
	AccessTokenLoader AccessTokenLoader
}

// NewClient creates a new instance of Client
func NewClient(options *ClientOptions) Client {
	if options == nil {
		options = &ClientOptions{}
	}

	return &helixClient{
		accessTokenLoader: options.AccessTokenLoader,
		clientID:          options.ClientID,
		Client: client.NewClient(&client.Options{
			RequestTimeout: options.RequestTimeout,
			Transport:      options.Transport,
		}),
	}
}

func (c *helixClient) headers(options *RequestOptions) func(context.Context) (http.Header, error) {
	return func(ctx context.Context) (http.Header, error) {
		h := http.Header{}

		h.Set("Accept", "application/json")
		h.Set("Client-ID", c.clientID)

		if options != nil && options.Token != "" {
			setToken(h, options.Token)
		} else if c.accessTokenLoader != nil {
			token, err := c.accessTokenLoader(ctx)
			if err != nil {
				return nil, err
			}
			if token != "" {
				setToken(h, token)
			}
		}

		return h, nil
	}
}

func setToken(h http.Header, token string) {
	h.Set("Authorization", "Bearer "+token)
}
