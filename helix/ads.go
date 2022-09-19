package helix

import (
	"context"
	"net/http"

	"github.com/aidenwallis/go-twitch-client/internal/client"
)

// Ads defines the ads namespace
type Ads interface {
	// StartCommercial implements https://dev.twitch.tv/docs/api/reference#start-commercial
	StartCommercial(context.Context, *StartCommercialRequest) (*StartCommercialResponse, error)
}

const commercialPath = "https://api.twitch.tv/helix/channels/commercial"

// StartCommercialRequest defines the options passed to StartCommercial
type StartCommercialRequest struct {
	*RequestOptions

	// BroadcasterID is the ID of the channel requesting a commercial
	BroadcasterID string

	// Length is the Desired length of the commercial in seconds. Valid options are 30, 60, 90, 120, 150, 180.
	Length int
}

// StartCommercialResponse defines the API response returned by StartCommercial
type StartCommercialResponse struct {
	// Data is a slice of Commercial
	Data []*Commercial `json:"data"`
}

// startCommercialBody implements the request body structure for StartCommercial
type startCommercialBody struct {
	BroadcasterID string `json:"broadcaster_id"`
	Length        int    `json:"length"`
}

// Commercial represents a commercial entity in Helix
type Commercial struct {
	// Length of the triggered commercial
	Length int `json:"length"`

	// RetryAfter is the Seconds until the next commercial can be served on this channel
	RetryAfter int `json:"retry_after"`

	// Message provides contextual information on why the request failed
	Message string `json:"message"`
}

// StartCommercial implements https://dev.twitch.tv/docs/api/reference#start-commercial
func (c *helixClient) StartCommercial(ctx context.Context, req *StartCommercialRequest) (*StartCommercialResponse, error) {
	return client.WithBody[StartCommercialResponse](c.Request(&client.RequestConfig{
		Method:  http.MethodPost,
		URL:     commercialPath,
		Headers: c.headers(req.RequestOptions),
	}).BodyJSON(&startCommercialBody{
		BroadcasterID: req.BroadcasterID,
		Length:        req.Length,
	}).Do(ctx))
}
