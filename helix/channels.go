package helix

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/aidenwallis/go-twitch-client/internal/client"
)

// Channels represents the channels namespace
type Channels interface {
	// GetChannelInformation implements https://dev.twitch.tv/docs/api/reference#get-channel-information
	//
	// Gets channel information for users.
	GetChannelInformation(context.Context, *GetChannelInformationRequest) (*GetChannelInformationResponse, error)

	// ModifyChannelInformation implements https://dev.twitch.tv/docs/api/reference#modify-channel-information
	//
	// Modifies channel information for users.
	ModifyChannelInformation(context.Context, *ModifyChannelInformationRequest) error

	// GetChannelEditors implements https://dev.twitch.tv/docs/api/reference#get-channel-editors
	//
	// Gets a list of users who have editor permissions for a specific channel.
	GetChannelEditors(context.Context, *GetChannelEditorsRequest) (*GetChannelEditorsResponse, error)
}

const channelsPath = "https://api.twitch.tv/helix/channels"

// GetChannelInformationRequest defines the options passed to GetChannelInformation
type GetChannelInformationRequest struct {
	*RequestOptions

	// BroadcasterIDs is the The ID of the broadcaster whose channel you want to get. You may specify a maximum of 100 IDs.
	BroadcasterIDs []string
}

// GetChannelInformationResponse represents the API response returned by GetChannelInformation
type GetChannelInformationResponse struct {
	// Data represents a slice of Channel
	Data []*Channel `json:"data"`
}

// Channel represents a helix Channel entity
type Channel struct {
	// BroadcasterID is the Twitch User ID of this channel owner.
	BroadcasterID string `json:"broadcaster_id"`

	// BroadcasterLogin is the Broadcaster’s user login name.
	BroadcasterLogin string `json:"broadcaster_login"`

	// BroadcasterName is the Twitch user display name of this channel owner.
	BroadcasterName string `json:"broadcaster_name"`

	// GameName is the Name of the game being played on the channel.
	GameName string `json:"game_name"`

	// GameID is the Current game ID being played on the channel.
	GameID string `json:"game_id"`

	// BroadcasterLanguage is the Language of the channel. A language value is either the ISO 639-1 two-letter code for
	// a supported stream language or “other”.
	BroadcasterLanguage string `json:"broadcaster_language"`

	// Title of the stream.
	Title string `json:"title"`

	// Delay is the Stream delay in seconds.
	Delay int `json:"delay"`
}

// GetChannelInformation implements https://dev.twitch.tv/docs/api/reference#get-channel-information
//
// Gets channel information for users.
func (c *helixClient) GetChannelInformation(ctx context.Context, req *GetChannelInformationRequest) (*GetChannelInformationResponse, error) {
	values := url.Values{}
	values["broadcaster_id"] = req.BroadcasterIDs

	return client.WithBody[GetChannelInformationResponse](c.Request(&client.RequestConfig{
		Method:  http.MethodGet,
		URL:     channelsPath,
		Headers: c.headers(req.RequestOptions),
		Query:   values,
	}).Do(ctx))
}

// ModifyChannelInformationRequest defines the options passed to ModifyChannelInformation
type ModifyChannelInformationRequest struct {
	*RequestOptions

	// BroadcasterID is the ID of the channel to be updated
	BroadcasterID string

	// GameID is The current game ID being played on the channel. Use “0” or “” (an empty string) to unset the game.
	GameID *string

	// BroadcasterLanguage is The language of the channel. A language value must be either the ISO 639-1 two-letter code for
	// a supported stream language or “other”.
	BroadcasterLanguage *string

	// Title is The title of the stream. Value must not be an empty string.
	Title *string

	// Delay is the Stream delay in seconds. Stream delay is a Twitch Partner feature; trying to set this value for other
	// account types will return a 400 error.
	Delay *int
}

// modifyChannelInformationBody implements the request body structure for PATCH channelsPath
type modifyChannelInformationBody struct {
	GameID              *string `json:"game_id,omitempty"`
	BroadcasterLanguage *string `json:"broadcaster_language,omitempty"`
	Title               *string `json:"title,omitempty"`
	Delay               *int    `json:"delay,omitempty"`
}

// ModifyChannelInformation implements https://dev.twitch.tv/docs/api/reference#modify-channel-information
//
// Modifies channel information for users.
func (c *helixClient) ModifyChannelInformation(ctx context.Context, req *ModifyChannelInformationRequest) error {
	values := url.Values{}
	values.Set("broadcaster_id", req.BroadcasterID)

	return client.WithoutBody(c.Request(&client.RequestConfig{
		Method:  http.MethodPatch,
		URL:     channelsPath,
		Headers: c.headers(req.RequestOptions),
		Query:   values,
	}).BodyJSON(&modifyChannelInformationBody{
		GameID:              req.GameID,
		BroadcasterLanguage: req.BroadcasterLanguage,
		Title:               req.Title,
		Delay:               req.Delay,
	}).Do(ctx))
}

const channelEditorsPath = "https://api.twitch.tv/helix/channels/editors"

// GetChannelEditorsRequest defines the options passed to GetChannelEditors
type GetChannelEditorsRequest struct {
	*RequestOptions

	// BroadcasterID is the Broadcaster’s user ID associated with the channel.
	BroadcasterID string
}

// GetChannelEditorsResponse defines the API response returned by GetChannelEditors
type GetChannelEditorsResponse struct {
	// Data represents a slice of channel editors
	Data []*ChannelEditor `json:"data"`
}

// ChannelEditor represents a channel editor in Helix
type ChannelEditor struct {
	// UserID is the User ID of the editor.
	UserID string `json:"user_id"`

	// UserName is the Display name of the editor.
	UserName string `json:"user_name"`

	// CreatedAt is the Date and time the editor was given editor permissions.
	CreatedAt time.Time `json:"created_at"`
}

// GetChannelEditors implements https://dev.twitch.tv/docs/api/reference#get-channel-editors
//
// Gets a list of users who have editor permissions for a specific channel.
func (c *helixClient) GetChannelEditors(ctx context.Context, req *GetChannelEditorsRequest) (*GetChannelEditorsResponse, error) {
	values := url.Values{}
	values.Set("broadcaster_id", req.BroadcasterID)

	return client.WithBody[GetChannelEditorsResponse](c.Request(&client.RequestConfig{
		Method:  http.MethodGet,
		URL:     channelEditorsPath,
		Headers: c.headers(req.RequestOptions),
		Query:   values,
	}).Do(ctx))
}
