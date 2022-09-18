package helix

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aidenwallis/go-twitch-client/internal/client"
)

// Users defines the users namespace
type Users interface {
	// BlockUser implements https://dev.twitch.tv/docs/api/reference#block-user
	BlockUser(ctx context.Context, req *BlockUserRequest) error

	// GetUserBlocks implements https://dev.twitch.tv/docs/api/reference#get-user-block-list
	GetUserBlocks(context.Context, *GetUserBlocksRequest) (*GetUserBlocksResponse, error)

	// GetUserFollows implements https://dev.twitch.tv/docs/api/reference#get-users-follows
	GetUserFollows(context.Context, *GetUserFollowsRequest) (*GetUserFollowsResponse, error)

	// GetUsers implements https://dev.twitch.tv/docs/api/reference#get-users
	GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error)

	// UnblockUser implements https://dev.twitch.tv/docs/api/reference#unblock-user
	UnblockUser(ctx context.Context, req *UnblockUserRequest) error
}

const usersPath = "https://api.twitch.tv/helix/users"

// GetUsersRequest is the set of options passed to GetUsers
type GetUsersRequest struct {
	*RequestOptions

	// IDs is the User ID. Multiple user IDs can be specified. Limit: 100.
	IDs []string

	// Logins is the User login name. Multiple login names can be specified. Limit: 100.
	Logins []string
}

// GetUsersResponse is the response returned by GetUsers
type GetUsersResponse struct {
	Data []*User `json:"data"`
}

// User implements the Helix user type
type User struct {
	// ID is the User’s ID.
	ID string `json:"id"`

	// Login is the User’s login name.
	Login string `json:"login"`

	// DisplayName is the User’s display name.
	DisplayName string `json:"display_name"`

	// Type is the User’s type: "staff", "admin", "global_mod", or ""
	Type string `json:"type"`

	// BroadcasterType is the User’s broadcaster type: "partner", "affiliate", or "".
	BroadcasterType string `json:"broadcaster_type"`

	// Description is the User’s channel description.
	Description string `json:"description"`

	// ProfileImageURL is the URL of the user’s profile image.
	ProfileImageURL string `json:"profile_image_url"`

	// OfflineImageURL is the URL of the user’s offline image.
	OfflineImageURL string `json:"offline_image_url"`

	// Email is the User’s verified email address. Returned if the request includes the user:read:email scope.
	Email string `json:"email,omitempty"`

	// CreatedAt is the Date when the user was created.
	CreatedAt time.Time `json:"created_at"`
}

// GetUsers implements https://dev.twitch.tv/docs/api/reference#get-users
func (c *helixClient) GetUsers(ctx context.Context, req *GetUsersRequest) (*GetUsersResponse, error) {
	values := url.Values{}
	for _, v := range req.IDs {
		values.Add("id", v)
	}
	for _, v := range req.Logins {
		values.Add("login", v)
	}

	return client.WithBody[GetUsersResponse](c.Request(&client.RequestConfig{
		Headers: c.headers(req.RequestOptions),
		Method:  http.MethodGet,
		URL:     usersPath,
		Query:   values,
	}).Do(ctx))
}

const userFollowsPath = "https://api.twitch.tv/helix/users/follows"

// GetUserFollowsRequest is the set of options passed to GetUserFollows
type GetUserFollowsRequest struct {
	*RequestOptions

	// After is a Cursor for forward pagination: tells the server where to start fetching the next set of results, in a multi-page response.
	// The cursor value specified here is from the pagination response field of a prior query.
	After string

	// FromID is a User ID. The request returns information about users who are being followed by the from_id user.
	FromID string

	// ToID is a User ID. The request returns information about users who are following the to_id user.
	ToID string

	// First is the Maximum number of objects to return. Maximum: 100. Default: 20.
	First int
}

// GetUserFollowsResponse is the response returned by GetUserFollows
type GetUserFollowsResponse struct {
	// Total is the number of items returned.
	//
	// * If only from_id was in the request, this is the total number of followed users.
	//
	// * If only to_id was in the request, this is the total number of followers.
	//
	// * If both from_id and to_id were in the request, this is 1 (if the "from" user follows the "to" user) or 0.
	Total int `json:"total"`

	// Data is a slice of UserFollow
	Data []*UserFollow `json:"data"`

	// Pagination contains a cursor value, to be used in a subsequent request to specify the starting point of the next set of results.
	Pagination Pagination `json:"pagination"`
}

// UserFollow represents a Helix user follow
type UserFollow struct {
	// FromID is the ID of the user following the to_id user.
	FromID string `json:"from_id"`

	// FromLogin is the Login of the user following the to_id user.
	FromLogin string `json:"from_login"`

	// FromName is the Display name corresponding to from_id.
	FromName string `json:"from_name"`

	// ToID is the ID of the user being followed by the from_id user.
	ToID string `json:"to_id"`

	// ToLogin is the Login of the user being followed by the from_id user.
	ToLogin string `json:"to_login"`

	// ToName is the Display name corresponding to to_id.
	ToName string `json:"to_name"`

	// FollowedAt is the Date and time when the from_id user followed the to_id user.
	FollowedAt time.Time `json:"followed_at"`
}

// GetUserFollows implements https://dev.twitch.tv/docs/api/reference#get-users-follows
func (c *helixClient) GetUserFollows(ctx context.Context, req *GetUserFollowsRequest) (*GetUserFollowsResponse, error) {
	values := url.Values{}
	if req.After != "" {
		values.Set("after", req.After)
	}
	if req.First > 0 {
		values.Set("first", strconv.Itoa(req.First))
	}
	if req.FromID != "" {
		values.Set("from_id", req.FromID)
	}
	if req.ToID != "" {
		values.Set("to_id", req.ToID)
	}

	return client.WithBody[GetUserFollowsResponse](c.Request(&client.RequestConfig{
		Headers: c.headers(req.RequestOptions),
		Method:  http.MethodGet,
		URL:     userFollowsPath,
		Query:   values,
	}).Do(ctx))
}

const userBlocksPath = "https://api.twitch.tv/helix/users/blocks"

// GetUserBlocksRequest is the set of options passed to GetUserBlocks
type GetUserBlocksRequest struct {
	*RequestOptions

	// BroadcasterID is the User ID for a Twitch user. This must match your access tokens' user.
	BroadcasterID string

	// After (optional) is a Cursor for forward pagination: tells the server where to start fetching the next set of results, in a multi-page response.
	// The cursor value specified here is from the pagination response field of a prior query.
	After string

	// First (optional) is the Maximum number of objects to return. Maximum: 100. Default: 20.
	First int
}

// GetUserBlocksResponse is the response returned by GetUserBlocks
type GetUserBlocksResponse struct {
	// Total is the number of items returned.
	//
	// * If only from_id was in the request, this is the total number of followed users.
	//
	// * If only to_id was in the request, this is the total number of followers.
	//
	// * If both from_id and to_id were in the request, this is 1 (if the "from" user follows the "to" user) or 0.
	Total int `json:"total"`

	// Data is a slice of UserBlock
	Data []*UserBlock `json:"data"`

	// Pagination contains a cursor value, to be used in a subsequent request to specify the starting point of the next set of results.
	Pagination Pagination `json:"pagination"`
}

// UserBlock represents a Helix user block
type UserBlock struct {
	// UserID is the User ID of the blocked user.
	UserID string `json:"user_id"`

	// UserLogin is the Login of the blocked user.
	UserLogin string `json:"user_login"`

	// DisplayName is the Display name of the blocked user.
	DisplayName string `json:"display_name"`
}

// GetUserBlocks implements https://dev.twitch.tv/docs/api/reference#get-user-block-list
func (c *helixClient) GetUserBlocks(ctx context.Context, req *GetUserBlocksRequest) (*GetUserBlocksResponse, error) {
	values := url.Values{}
	values.Set("broadcaster_id", req.BroadcasterID)
	if req.After != "" {
		values.Set("after", req.After)
	}
	if req.First > 0 {
		values.Set("first", strconv.Itoa(req.First))
	}

	return client.WithBody[GetUserBlocksResponse](c.Request(&client.RequestConfig{
		Headers: c.headers(req.RequestOptions),
		Method:  http.MethodGet,
		URL:     userBlocksPath,
		Query:   values,
	}).Do(ctx))
}

// BlockUserRequest is the set of options passed to BlockUser
type BlockUserRequest struct {
	*RequestOptions

	// TargetUserID is the User ID of the user to be blocked.
	TargetUserID string

	// SourceContext (optional) is the Source context for blocking the user. Valid values: "chat", "whisper".
	SourceContext string

	// Reason (optional) is the Reason for blocking the user. Valid values: "spam", "harassment", or "other".
	Reason string
}

// BlockUser implements https://dev.twitch.tv/docs/api/reference#block-user
func (c *helixClient) BlockUser(ctx context.Context, req *BlockUserRequest) error {
	values := url.Values{}
	values.Set("target_user_id", req.TargetUserID)
	if req.SourceContext != "" {
		values.Set("source_context", req.SourceContext)
	}
	if req.Reason != "" {
		values.Set("reason", req.Reason)
	}

	return client.WithoutBody(c.Request(&client.RequestConfig{
		Headers: c.headers(req.RequestOptions),
		Method:  http.MethodPut,
		URL:     userBlocksPath,
		Query:   values,
	}).Do(ctx))
}

// UnblockUserRequest is the set of options passed to UnblockUser
type UnblockUserRequest struct {
	*RequestOptions

	// TargetUserID is the User ID of the user to be blocked.
	TargetUserID string
}

// UnblockUser implements https://dev.twitch.tv/docs/api/reference#unblock-user
func (c *helixClient) UnblockUser(ctx context.Context, req *UnblockUserRequest) error {
	values := url.Values{}
	values.Set("target_user_id", req.TargetUserID)

	return client.WithoutBody(c.Request(&client.RequestConfig{
		Headers: c.headers(req.RequestOptions),
		Method:  http.MethodDelete,
		URL:     userBlocksPath,
		Query:   values,
	}).Do(ctx))
}
