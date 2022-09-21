package helix

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/aidenwallis/go-twitch-client"
	"github.com/aidenwallis/go-twitch-client/internal/client"
)

// Chat represents the chat namespace
type Chat interface {
	// GetChannelEmotes implements https://dev.twitch.tv/docs/api/reference#get-channel-emotes
	//
	// Gets all emotes that the specified Twitch channel created. Broadcasters create these custom emotes for users who subscribe to or
	// follow the channel, or cheer Bits in the channel’s chat window.
	GetChannelEmotes(context.Context, *GetChannelEmotesRequest) (*GetChannelEmotesResponse, error)

	// GetGlobalEmotes implements https://dev.twitch.tv/docs/api/reference#get-global-emotes
	//
	// Gets all global emotes. Global emotes are Twitch-created emoticons that users can use in any Twitch chat.
	GetGlobalEmotes(context.Context, *GetGlobalEmotesRequest) (*GetGlobalEmotesResponse, error)

	// GetEmoteSets implements https://dev.twitch.tv/docs/api/reference#get-emote-sets
	//
	// Gets emotes for one or more specified emote sets.
	//
	// An emote set groups emotes that have a similar context. For example, Twitch places all the subscriber emotes that a broadcaster
	// uploads for their channel in the same emote set.
	//
	// Learn more: https://dev.twitch.tv/docs/irc/emotes
	GetEmoteSets(context.Context, *GetEmoteSetsRequest) (*GetEmoteSetsResponse, error)

	// GetChannelChatBadges implements https://dev.twitch.tv/docs/api/reference#get-channel-chat-badges
	//
	// Gets a list of custom chat badges that can be used in chat for the specified channel. This includes subscriber badges and Bit badges.
	GetChannelChatBadges(context.Context, *GetChannelChatBadgesRequest) (*GetChannelChatBadgesResponse, error)

	// GetChannelChatBadges implements https://dev.twitch.tv/docs/api/reference#get-channel-chat-badges
	//
	// Gets a list of custom chat badges that can be used in chat for the specified channel. This includes subscriber badges and Bit badges.
	GetGlobalChatBadges(context.Context, *GetGlobalChatBadgesRequest) (*GetGlobalChatBadgesResponse, error)

	// GetChatSettings implements https://dev.twitch.tv/docs/api/reference#get-chat-settings
	//
	// Gets the broadcaster’s chat settings.
	//
	// For an overview of chat settings, see:
	//
	// * Chat Commands for Broadcasters and Moderators: https://help.twitch.tv/s/article/chat-commands#AllMods
	//
	// * Moderator Preferences: https://help.twitch.tv/s/article/setting-up-moderation-for-your-twitch-channel#modpreferences
	GetChatSettings(context.Context, *GetChatSettingsRequest) (*GetChatSettingsResponse, error)

	// UpdateChatSettings implements https://dev.twitch.tv/docs/api/reference#update-chat-settings
	//
	// Updates the broadcaster’s chat settings.
	UpdateChatSettings(context.Context, *UpdateChatSettingsRequest) (*UpdateChatSettingsResponse, error)

	// SendChatAnnouncement implements https://dev.twitch.tv/docs/api/reference#send-chat-announcement
	//
	// Sends an announcement to the broadcaster’s chat room.
	SendChatAnnouncement(context.Context, *SendChatAnnouncementRequest) error

	// GetUserChatColors implements https://dev.twitch.tv/docs/api/reference#get-user-chat-color
	//
	// Gets the color used for the user’s name in chat.
	GetUserChatColors(context.Context, *GetUserChatColorsRequest) (*GetUserChatColorsResponse, error)

	// UpdateUserChatColor implements https://dev.twitch.tv/docs/api/reference#update-user-chat-color
	//
	// Updates the color used for the user’s name in chat.
	UpdateUserChatColor(context.Context, *UpdateUserChatColorRequest) error
}

const chatEmotesPath = "https://api.twitch.tv/helix/chat/emotes"

// GetChannelEmotesRequest defines the options passed to GetChannelEmotes
type GetChannelEmotesRequest struct {
	*RequestOptions

	// BroadcasterID is an ID that identifies the broadcaster to get the emotes from.
	BroadcasterID string
}

// GetChannelEmotesResponse defines the API response returned by GetChannelEmotes
type GetChannelEmotesResponse struct {
	// Data defines a slice of ChatEmote
	Data []*ChannelEmote `json:"data"`

	// Template is a templated URL. Use the values from id, format, scale, and theme_mode to replace the like-named placeholder strings
	// in the templated URL to create a CDN (content delivery network) URL that you use to fetch the emote. For information about what
	// the template looks like and how to use it to fetch emotes, see: https://dev.twitch.tv/docs/irc/emotes#cdn-template
	//
	// You may use EmoteURL() on this struct to help you generate a URL.
	Template string `json:"template"`
}

// FormatChatEmoteTemplate returns a CDN emote URL based off of a defined template and emote properties
func FormatChatEmoteTemplate(template, id, format, themeMode, scale string) string {
	return strings.NewReplacer(
		"{{id}}", id,
		"{{format}}", format,
		"{{theme_mode}}", themeMode,
		"{{scale}}", scale,
	).Replace(template)
}

// EmoteURL returns a valid emote CDN URL based on the template returned in the Helix response.
func (r *GetChannelEmotesResponse) EmoteURL(id, format, themeMode, scale string) string {
	return FormatChatEmoteTemplate(r.Template, id, format, themeMode, scale)
}

// ChannelEmote defines a single channel chat emote in Helix
type ChannelEmote struct {
	// ID that identifies the emote.
	ID string `json:"id"`

	// Name of the emote. This is the name that viewers type in the chat window to get the emote to appear.
	Name string `json:"name"`

	// Tier is the subscriber tier at which the emote is unlocked. This field contains the tier information only if emote_type is set to
	// subscriptions, otherwise, it’s an empty string.
	Tier string `json:"tier"`

	// EmoteType is the type of emote. The possible values are:
	//
	// * bitstier — Indicates a custom Bits tier emote.
	//
	// * follower — Indicates a custom follower emote.
	//
	// * subscriptions — Indicates a custom subscriber emote.
	EmoteType string `json:"emote_type"`

	// EmoteSetID is an ID that identifies the emote set that the emote belongs to.
	EmoteSetID string `json:"emote_set_id"`

	// Format is the the formats that the emote is available in. For example, if the emote is available only as a static PNG, the array
	// contains only static. But if it’s available as a static PNG and an animated GIF, the array contains static and animated. The possible
	// formats are:
	//
	// * animated — Indicates an animated GIF is available for this emote.
	//
	// * static — Indicates a static PNG file is available for this emote.
	Format []string `json:"format"`

	// Scale is the sizes that the emote is available in. For example, if the emote is available in small and medium sizes, the array contains
	// 1.0 and 2.0. Possible sizes are:
	//
	// * 1.0 — A small version (28px x 28px) is available.
	//
	// * 2.0 — A medium version (56px x 56px) is available.
	//
	// * 3.0 — A large version (112px x 112px) is available.
	Scale []string `json:"scale"`

	// ThemeMode is the background themes that the emote is available in. Possible themes are:
	//
	// * dark
	//
	// * light
	ThemeMode []string `json:"theme_mode"`

	// Images contains the image URLs for the emote. These image URLs will always provide a static (i.e., non-animated) emote image with a
	// light background.
	//
	// Deprecated: **NOTE** The preference is for you to use the templated URL in the template field to fetch the image instead of using these URLs.
	Images ChatEmoteImages `json:"images"`
}

// ChatEmoteImages defines the emote CDN URLs for a given ChatEmote
//
// Deprecated: This is deprecated in favor of the URL template returned by Helix
type ChatEmoteImages struct {
	// URL1x is a URL to the small version (28px x 28px) of the emote.
	URL1x string `json:"url_1x"`

	// URL1x is a URL to the medium version (56px x 56px) of the emote.
	URL2x string `json:"url_2x"`

	// URL1x is a URL to the large version (112px x 112px) of the emote.
	URL3x string `json:"url_3x"`
}

// GetChannelEmotes implements https://dev.twitch.tv/docs/api/reference#get-channel-emotes
//
// Gets all emotes that the specified Twitch channel created. Broadcasters create these custom emotes for users who subscribe to or
// follow the channel, or cheer Bits in the channel’s chat window.
func (c *helixClient) GetChannelEmotes(ctx context.Context, req *GetChannelEmotesRequest) (*GetChannelEmotesResponse, error) {
	values := url.Values{}
	values.Set("broadcaster_id", req.BroadcasterID)

	return client.WithBody[GetChannelEmotesResponse](c.Request(&client.RequestConfig{
		Method:  http.MethodGet,
		URL:     chatEmotesPath,
		Headers: c.headers(req.RequestOptions),
		Query:   values,
	}).Do(ctx))
}

const chatGlobalEmotesPath = "https://api.twitch.tv/helix/chat/emotes/global"

// GetGlobalEmotesRequest defines the options passed to GetGlobalEmotes
type GetGlobalEmotesRequest struct {
	*RequestOptions
}

// GetChannelEmotesResponse defines the API response returned by GetChannelEmotes
type GetGlobalEmotesResponse struct {
	// Data defines a slice of ChatEmote
	Data []*GlobalEmote `json:"data"`

	// Template is a templated URL. Use the values from id, format, scale, and theme_mode to replace the like-named placeholder strings
	// in the templated URL to create a CDN (content delivery network) URL that you use to fetch the emote. For information about what
	// the template looks like and how to use it to fetch emotes, see: https://dev.twitch.tv/docs/irc/emotes#cdn-template
	//
	// You may use EmoteURL() on this struct to help you generate a URL.
	Template string `json:"template"`
}

// EmoteURL returns a valid emote CDN URL based on the template returned in the Helix response.
func (r *GetGlobalEmotesResponse) EmoteURL(id, format, themeMode, scale string) string {
	return FormatChatEmoteTemplate(r.Template, id, format, themeMode, scale)
}

// GlobalEmote defines a single global emote in Helix
type GlobalEmote struct {
	// ID that identifies the emote.
	ID string `json:"id"`

	// Name of the emote. This is the name that viewers type in the chat window to get the emote to appear.
	Name string `json:"name"`

	// Format is the the formats that the emote is available in. For example, if the emote is available only as a static PNG, the array
	// contains only static. But if it’s available as a static PNG and an animated GIF, the array contains static and animated. The possible
	// formats are:
	//
	// * animated — Indicates an animated GIF is available for this emote.
	//
	// * static — Indicates a static PNG file is available for this emote.
	Format []string `json:"format"`

	// Scale is the sizes that the emote is available in. For example, if the emote is available in small and medium sizes, the array contains
	// 1.0 and 2.0. Possible sizes are:
	//
	// * 1.0 — A small version (28px x 28px) is available.
	//
	// * 2.0 — A medium version (56px x 56px) is available.
	//
	// * 3.0 — A large version (112px x 112px) is available.
	Scale []string `json:"scale"`

	// ThemeMode is the background themes that the emote is available in. Possible themes are:
	//
	// * dark
	//
	// * light
	ThemeMode []string `json:"theme_mode"`

	// Images contains the image URLs for the emote. These image URLs will always provide a static (i.e., non-animated) emote image with a
	// light background.
	//
	// Deprecated: **NOTE** The preference is for you to use the templated URL in the template field to fetch the image instead of using these URLs.
	Images ChatEmoteImages `json:"images"`
}

// GetGlobalEmotes implements https://dev.twitch.tv/docs/api/reference#get-global-emotes
//
// Gets all global emotes. Global emotes are Twitch-created emoticons that users can use in any Twitch chat.
func (c *helixClient) GetGlobalEmotes(ctx context.Context, req *GetGlobalEmotesRequest) (*GetGlobalEmotesResponse, error) {
	return client.WithBody[GetGlobalEmotesResponse](c.Request(&client.RequestConfig{
		Method:  http.MethodGet,
		URL:     chatGlobalEmotesPath,
		Headers: c.headers(req.RequestOptions),
	}).Do(ctx))
}

const chatEmoteSetsPath = "https://api.twitch.tv/helix/chat/emotes/set"

// GetEmoteSetsRequest defines the options passed to GetEmoteSets
type GetEmoteSetsRequest struct {
	*RequestOptions

	// EmoteSetIDs are IDs that identify the emote set. Include the parameter for each emote set you want to get.
	// You may specify a maximum of 25 IDs.
	EmoteSetIDs []string
}

// GetEmoteSetsResponse defines the API response returned by GetEmoteSets
type GetEmoteSetsResponse struct {
	// Data represents a slice of Emote Sets
	Data []*SetEmote `json:"data"`

	// Template is a templated URL. Use the values from id, format, scale, and theme_mode to replace the like-named placeholder strings
	// in the templated URL to create a CDN (content delivery network) URL that you use to fetch the emote. For information about what
	// the template looks like and how to use it to fetch emotes, see: https://dev.twitch.tv/docs/irc/emotes#cdn-template
	//
	// You may use EmoteURL() on this struct to help you generate a URL.
	Template string `json:"template"`
}

// EmoteURL returns a valid emote CDN URL based on the template returned in the Helix response.
func (r *GetEmoteSetsResponse) EmoteURL(id, format, themeMode, scale string) string {
	return FormatChatEmoteTemplate(r.Template, id, format, themeMode, scale)
}

// SetEmote defines a single emoteset emote in Helix
type SetEmote struct {
	// ID that identifies the emote.
	ID string `json:"id"`

	// Name of the emote. This is the name that viewers type in the chat window to get the emote to appear.
	Name string `json:"name"`

	// EmoteType is the type of emote. The possible values are:
	//
	// * bitstier — Indicates a custom Bits tier emote.
	//
	// * follower — Indicates a custom follower emote.
	//
	// * subscriptions — Indicates a custom subscriber emote.
	EmoteType string `json:"emote_type"`

	// EmoteSetID is an ID that identifies the emote set that the emote belongs to.
	EmoteSetID string `json:"emote_set_id"`

	// OwnerID is the ID of the broadcaster who owns the emote.
	OwnerID string `json:"owner_id"`

	// Format is the the formats that the emote is available in. For example, if the emote is available only as a static PNG, the array
	// contains only static. But if it’s available as a static PNG and an animated GIF, the array contains static and animated. The possible
	// formats are:
	//
	// * animated — Indicates an animated GIF is available for this emote.
	//
	// * static — Indicates a static PNG file is available for this emote.
	Format []string `json:"format"`

	// Scale is the sizes that the emote is available in. For example, if the emote is available in small and medium sizes, the array contains
	// 1.0 and 2.0. Possible sizes are:
	//
	// * 1.0 — A small version (28px x 28px) is available.
	//
	// * 2.0 — A medium version (56px x 56px) is available.
	//
	// * 3.0 — A large version (112px x 112px) is available.
	Scale []string `json:"scale"`

	// ThemeMode is the background themes that the emote is available in. Possible themes are:
	//
	// * dark
	//
	// * light
	ThemeMode []string `json:"theme_mode"`

	// Images contains the image URLs for the emote. These image URLs will always provide a static (i.e., non-animated) emote image with a
	// light background.
	//
	// Deprecated: **NOTE** The preference is for you to use the templated URL in the template field to fetch the image instead of using these URLs.
	Images ChatEmoteImages `json:"images"`
}

// GetEmoteSets implements https://dev.twitch.tv/docs/api/reference#get-emote-sets
//
// Gets emotes for one or more specified emote sets.
//
// An emote set groups emotes that have a similar context. For example, Twitch places all the subscriber emotes that a broadcaster uploads for their
// channel in the same emote set.
//
// Learn more: https://dev.twitch.tv/docs/irc/emotes
func (c *helixClient) GetEmoteSets(ctx context.Context, req *GetEmoteSetsRequest) (*GetEmoteSetsResponse, error) {
	values := url.Values{}
	values["emote_set_id"] = req.EmoteSetIDs

	return client.WithBody[GetEmoteSetsResponse](c.Request(&client.RequestConfig{
		Method:  http.MethodGet,
		URL:     chatEmoteSetsPath,
		Headers: c.headers(req.RequestOptions),
		Query:   values,
	}).Do(ctx))
}

const channelChatBadges = "https://api.twitch.tv/helix/chat/badges"

// GetChannelChatBadgesRequest defines the options passed to GetChannelChatBadges
type GetChannelChatBadgesRequest struct {
	*RequestOptions

	// BroadcasterID is the ID of the broadcaster whose chat badges you want to get.
	BroadcasterID string
}

// GetChannelChatBadgesResponse defines the API response returned by GetChannelChatBadges
type GetChannelChatBadgesResponse struct {
	// Data represents a slice of ChatBadge
	Data []*ChatBadge `json:"data"`
}

// ChatBadge defines a single helix chat badge
type ChatBadge struct {
	// SetID is the ID for the chat badge set.
	SetID string `json:"site_id"`

	// Versions contains chat badge objects for the set.
	Versions []*ChatBadgeVersion `json:"versions"`
}

// ChatBadgeVersion defines a single chat badge version in Helix
type ChatBadgeVersion struct {
	// ID of the chat badge version.
	ID string `json:"id"`

	// ImageURL1x is the Small image URL.
	ImageURL1x string `json:"image_url_1x"`

	// ImageURL2x is the Medium image URL.
	ImageURL2x string `json:"image_url_2x"`

	// ImageURL4x is the Large image URL.
	ImageURL4x string `json:"image_url_4x"`
}

// GetChannelChatBadges implements https://dev.twitch.tv/docs/api/reference#get-channel-chat-badges
//
// Gets a list of custom chat badges that can be used in chat for the specified channel. This includes subscriber badges and Bit badges.
func (c *helixClient) GetChannelChatBadges(ctx context.Context, req *GetChannelChatBadgesRequest) (*GetChannelChatBadgesResponse, error) {
	values := url.Values{}
	values.Set("broadcaster_id", req.BroadcasterID)

	return client.WithBody[GetChannelChatBadgesResponse](c.Request(&client.RequestConfig{
		Method:  http.MethodGet,
		URL:     channelChatBadges,
		Headers: c.headers(req.RequestOptions),
		Query:   values,
	}).Do(ctx))
}

const globalChatBadges = "https://api.twitch.tv/helix/chat/badges/global"

// GetGlobalChatBadgesRequest defines the options passed to GetGlobalChatBadges
type GetGlobalChatBadgesRequest struct {
	*RequestOptions
}

// GetGlobalChatBadgesResponse defines the API response returned by GetGlobalChatBadges
type GetGlobalChatBadgesResponse struct {
	// Data represents a slice of ChatBadge
	Data []*ChatBadge `json:"data"`
}

// GetChannelChatBadges implements https://dev.twitch.tv/docs/api/reference#get-channel-chat-badges
//
// Gets a list of custom chat badges that can be used in chat for the specified channel. This includes subscriber badges and Bit badges.
func (c *helixClient) GetGlobalChatBadges(ctx context.Context, req *GetGlobalChatBadgesRequest) (*GetGlobalChatBadgesResponse, error) {
	return client.WithBody[GetGlobalChatBadgesResponse](c.Request(&client.RequestConfig{
		Method:  http.MethodGet,
		URL:     globalChatBadges,
		Headers: c.headers(req.RequestOptions),
	}).Do(ctx))
}

const chatSettingsPath = "https://api.twitch.tv/helix/chat/settings"

// GetChatSettingsRequest defines the options passed to GetChatSettings
type GetChatSettingsRequest struct {
	*RequestOptions

	// BroadcasterID is the ID of the broadcaster whose chat settings you want to get.
	BroadcasterID string

	// ModeratorID (optional)
	//
	// Required only to access the non_moderator_chat_delay or non_moderator_chat_delay_duration settings.
	//
	// The ID of a user that has permission to moderate the broadcaster’s chat room. This ID must match the user ID associated with the user OAuth token.
	//
	// If the broadcaster wants to get their own settings (instead of having the moderator do it), set this parameter to the broadcaster’s ID, too.
	ModeratorID string
}

// GetChatSettingsResponse defines the API response returned by GetChatSettings
type GetChatSettingsResponse struct {
	// Data represents a slice of ChatSettings
	Data []*ChatSettings `json:"data"`
}

// ChatSettings defines the returned chat settings entity from Helix
type ChatSettings struct {
	// EmoteMode determines whether chat messages must contain only emotes.
	//
	// Is true, if only messages that are 100% emotes are allowed; otherwise, false.
	EmoteMode bool `json:"emote_mode"`

	// FollowerMode determines whether the broadcaster restricts the chat room to followers only, based on how long they’ve followed.
	//
	// Is true, if the broadcaster restricts the chat room to followers only; otherwise, false.
	//
	// See FollowerModeDuration for how long the followers must have followed the broadcaster to participate in the chat room.
	FollowerMode bool `json:"follower_mode"`

	// NonModeratorChatDelay that determines whether the broadcaster adds a short delay before chat messages appear in the chat room.
	// This gives chat moderators and bots a chance to remove them before viewers can see the message.
	//
	// Is true, if the broadcaster applies a delay; otherwise, false.
	//
	// See NonModeratorChatDelayDuration for the length of the delay.
	//
	// The response includes this field only if the request specifies a user access token that includes the moderator:read:chat_settings
	// scope and the user in the moderator_id query parameter is one of the broadcaster’s moderators.
	NonModeratorChatDelay bool `json:"non_moderator_chat_delay"`

	// SlowMode determines whether the broadcaster limits how often users in the chat room are allowed to send messages.
	//
	// Is true, if the broadcaster applies a delay; otherwise, false.
	//
	// See SlowModeWaitTime for the delay.
	SlowMode bool `json:"slow_mode"`

	// SubscriberMode determines whether only users that subscribe to the broadcaster’s channel can talk in the chat room.
	//
	// Is true, if the broadcaster restricts the chat room to subscribers only; otherwise, false.
	SubscriberMode bool `json:"subscriber_mode"`

	// UniqueChatMode determines whether the broadcaster requires users to post only unique messages in the chat room.
	//
	// Is true, if the broadcaster requires unique messages only; otherwise, false.
	UniqueChatMode bool `json:"unique_chat_mode"`

	// FollowerModeDuration is the length of time, in minutes, that the followers must have followed the broadcaster to
	// participate in the chat room. See FollowerMode.
	//
	// Is nil if FollowerMode is false.
	FollowerModeDuration *int `json:"follower_mode_duration"`

	// The amount of time, in seconds, that messages are delayed from appearing in chat. See NonModeratorChatDelay.
	//
	// Is nil if NonModeratorChatDelay is false.
	//
	// The response includes this field only if the request specifies a user access token that includes the moderator:read:chat_settings
	// scope and the user in the moderator_id query parameter is one of the broadcaster’s moderators.
	NonModeratorChatDelayDuration *int `json:"non_moderator_chat_delay_duration"`

	// The amount of time, in seconds, that users need to wait between sending messages. See SlowMode.
	//
	// Is nil if SlowMode is false.
	SlowModeWaitTime *int `json:"slow_mode_wait_time"`
}

// GetChatSettings implements https://dev.twitch.tv/docs/api/reference#get-chat-settings
//
// Gets the broadcaster’s chat settings.
//
// For an overview of chat settings, see:
//
// * Chat Commands for Broadcasters and Moderators: https://help.twitch.tv/s/article/chat-commands#AllMods
//
// * Moderator Preferences: https://help.twitch.tv/s/article/setting-up-moderation-for-your-twitch-channel#modpreferences
func (c *helixClient) GetChatSettings(ctx context.Context, req *GetChatSettingsRequest) (*GetChatSettingsResponse, error) {
	values := url.Values{}
	values.Set("broadcaster_id", req.BroadcasterID)
	if req.ModeratorID != "" {
		values.Set("moderator_id", req.ModeratorID)
	}

	return client.WithBody[GetChatSettingsResponse](c.Request(&client.RequestConfig{
		Method:  http.MethodGet,
		URL:     chatSettingsPath,
		Headers: c.headers(req.RequestOptions),
		Query:   values,
	}).Do(ctx))
}

// UpdateChatSettingsRequest defines the settings passed to UpdateChatSettings
type UpdateChatSettingsRequest struct {
	*RequestOptions

	// BroadcasterID is the ID of the broadcaster whose chat settings you want to update
	BroadcasterID string

	// ModeratorID is the D of a user that has permission to moderate the broadcaster’s chat room.
	// This ID must match the user ID associated with the user OAuth token.
	//
	// If the broadcaster is making the update, also set this parameter to the broadcaster’s ID.
	ModeratorID string

	// EmoteMode is a value that determines whether chat messages must contain only emotes.
	//
	// Set to true, if only messages that are 100% emotes are allowed; otherwise, false. Default is false.
	EmoteMode *bool

	// FollowerMode is a value that determines whether the broadcaster restricts the chat room to followers only,
	// based on how long they’ve followed.
	//
	// Set to true, if the broadcaster restricts the chat room to followers only; otherwise, false. Default is true.
	//
	// See FollowerModeDuration for how long the followers must have followed the broadcaster to participate in the chat room.
	FollowerMode *bool

	// NonModeratorChatDelay is a value that determines whether the broadcaster adds a short delay before chat messages appear
	// in the chat room. This gives chat moderators and bots a chance to remove them before viewers can see the message.
	//
	// Set to true, if the broadcaster applies a delay; otherwise, false. Default is false.
	//
	// See NonModeratorChatDelayDuration for the length of the delay.
	NonModeratorChatDelay *bool

	// SlowMode is a value that determines whether the broadcaster limits how often users in the chat room are allowed to send messages.
	//
	// Set to true, if the broadcaster applies a wait period messages; otherwise, false. Default is false.
	//
	// See SlowModeWaitTime for the delay.
	SlowMode *bool

	// SubscriberMode is a value that determines whether only users that subscribe to the broadcaster’s channel can talk in the chat room.
	//
	// Set to true, if the broadcaster restricts the chat room to subscribers only; otherwise, false. Default is false.
	SubscriberMode *bool

	// UniqueChatMode is a value that determines whether the broadcaster requires users to post only unique messages in the chat room.
	//
	// Set to true, if the broadcaster requires unique messages only; otherwise, false. Default is false.
	UniqueChatMode *bool

	// FollowerModeDuration is the length of time, in minutes, that the followers must have followed the broadcaster to participate in the
	// chat room (see FollowerMode).
	//
	// You may specify a value in the range: 0 (no restriction) through 129600 (3 months). The default is 0.
	FollowerModeDuration *int

	// NonModeratorChatDelayDuration is the amount of time, in seconds, that messages are delayed from appearing in chat.
	//
	// Possible values are:
	//
	// * 2 — 2 second delay (recommended)
	//
	// * 4 — 4 second delay
	//
	// * 6 — 6 second delay
	//
	// See NonModeratorChatDelay.
	NonModeratorChatDelayDuration *int

	// The amount of time, in seconds, that users need to wait between sending messages (see SlowMode).
	//
	// You may specify a value in the range: 3 (3 second delay) through 120 (2 minute delay). The default is 30 seconds.
	SlowModeWaitTime *int
}

// UpdateChatSettingsResponse defines the API response returned by UpdateChatSettings
type UpdateChatSettingsResponse struct {
	// Data represents a slice of ChatSettings
	Data []*ChatSettings `json:"data"`
}

// updateChatSettingsBody is the request body defined for UpdateChatSettings
type updateChatSettingsBody struct {
	EmoteMode                     *bool `json:"emote_mode,omitempty"`
	FollowerMode                  *bool `json:"follower_mode,omitempty"`
	NonModeratorChatDelay         *bool `json:"non_moderator_chat_delay,omitempty"`
	SlowMode                      *bool `json:"slow_mode,omitempty"`
	SubscriberMode                *bool `json:"subscriber_mode,omitempty"`
	UniqueChatMode                *bool `json:"unique_chat_mode,omitempty"`
	FollowerModeDuration          *int  `json:"follower_mode_duration,omitempty"`
	NonModeratorChatDelayDuration *int  `json:"non_moderator_chat_delay_duration,omitempty"`
	SlowModeWaitTime              *int  `json:"slow_mode_wait_time,omitempty"`
}

// UpdateChatSettings implements https://dev.twitch.tv/docs/api/reference#update-chat-settings
//
// Updates the broadcaster’s chat settings.
func (c *helixClient) UpdateChatSettings(ctx context.Context, req *UpdateChatSettingsRequest) (*UpdateChatSettingsResponse, error) {
	values := url.Values{}
	values.Set("broadcaster_id", req.BroadcasterID)
	values.Set("moderator_id", req.ModeratorID)

	return client.WithBody[UpdateChatSettingsResponse](c.Request(&client.RequestConfig{
		Method:  http.MethodPatch,
		URL:     chatSettingsPath,
		Query:   values,
		Headers: c.headers(req.RequestOptions),
	}).BodyJSON(&updateChatSettingsBody{
		EmoteMode:                     req.EmoteMode,
		FollowerMode:                  req.FollowerMode,
		NonModeratorChatDelay:         req.NonModeratorChatDelay,
		SlowMode:                      req.SlowMode,
		SubscriberMode:                req.SubscriberMode,
		UniqueChatMode:                req.UniqueChatMode,
		FollowerModeDuration:          req.FollowerModeDuration,
		NonModeratorChatDelayDuration: req.NonModeratorChatDelayDuration,
		SlowModeWaitTime:              req.SlowModeWaitTime,
	}).Do(ctx))
}

const chatAnnouncementsPath = "https://api.twitch.tv/helix/chat/announcements"

// SendChatAnnouncementRequest defines the options passed to SendChatAnnouncement
type SendChatAnnouncementRequest struct {
	*RequestOptions

	// BroadcasterID is the ID of the broadcaster that owns the chat room to send the announcement to
	BroadcasterID string

	// ModeratorID is the ID of a user who has permission to moderate the broadcaster’s chat room.
	// This ID must match the user ID in the OAuth token, which can be a moderator or the broadcaster.
	ModeratorID string

	// Color (optional) used to highlight the announcement. Possible case-sensitive values are:
	//
	// * blue
	//
	// * green
	//
	// * orange
	//
	// * purple
	//
	// * primary (default)
	//
	// If color is set to primary or is not set, the channel’s accent color is used to highlight the announcement (see
	// Profile Accent Color under profile settings, Channel and Videos, and Brand): https://twitch.tv/settings/profile
	Color *string

	// Message is the announcement to make in the broadcaster’s chat room. Announcements are limited to a maximum of 500 characters;
	// announcements longer than 500 characters are truncated.
	Message string
}

// sendChatAnnouncementBody defines the request body for SendChatAnnouncement
type sendChatAnnouncementBody struct {
	Color   string `json:"color,omitempty"`
	Message string `json:"message"`
}

// SendChatAnnouncement implements https://dev.twitch.tv/docs/api/reference#send-chat-announcement
//
// Sends an announcement to the broadcaster’s chat room.
func (c *helixClient) SendChatAnnouncement(ctx context.Context, req *SendChatAnnouncementRequest) error {
	values := url.Values{}
	values.Set("broadcaster_id", req.BroadcasterID)
	values.Set("moderator_id", req.ModeratorID)

	return client.WithoutBody(c.Request(&client.RequestConfig{
		Method:  http.MethodPost,
		URL:     chatAnnouncementsPath,
		Headers: c.headers(req.RequestOptions),
		Query:   values,
	}).BodyJSON(&sendChatAnnouncementBody{
		Color:   twitch.PointerValue(req.Color),
		Message: req.Message,
	}).Do(ctx))
}

const chatColorPath = "https://api.twitch.tv/helix/chat/color"

// GetUserChatColorsRequest defines the options passed to GetUserChatColors
type GetUserChatColorsRequest struct {
	*RequestOptions

	// UserIDs are the IDs of the users whose colors you want to get.
	// The maximum number of IDs that you may specify is 100.
	UserIDs []string
}

// UserChatColor defines a users' chat color in Helix
type UserChatColor struct {
	// UserID is the ID of the user
	UserID string `json:"user_id"`

	// UserLogin is the user’s login name.
	UserLogin string `json:"user_login"`

	// UserName is the user’s display name.
	UserName string `json:"user_name"`

	// Color is the Hex color code that the user uses in chat for their name.
	// If the user hasn’t specified a color in their settings, the string is empty.
	Color string `json:"color"`
}

// GetUserChatColorReponse defines the API response returned by GetUserChatColors
type GetUserChatColorsResponse struct {
	// Data represents a slice of UserChatColor
	Data []*UserChatColor `json:"data"`
}

// GetUserChatColors implements https://dev.twitch.tv/docs/api/reference#get-user-chat-color
//
// Gets the color used for the user’s name in chat.
func (c *helixClient) GetUserChatColors(ctx context.Context, req *GetUserChatColorsRequest) (*GetUserChatColorsResponse, error) {
	values := url.Values{}
	values["user_id"] = req.UserIDs

	return client.WithBody[GetUserChatColorsResponse](c.Request(&client.RequestConfig{
		Method:  http.MethodGet,
		URL:     chatColorPath,
		Headers: c.headers(req.RequestOptions),
		Query:   values,
	}).Do(ctx))
}

// UpdateUserChatColorRequest implements https://dev.twitch.tv/docs/api/reference#update-user-chat-color
//
//	Updates the color used for the user’s name in chat.
type UpdateUserChatColorRequest struct {
	*RequestOptions

	// UserID is te ID of the user whose chat color you want to update
	UserID string

	// Color to use for the user’s name in chat. All users may specify one of the following named color values.
	//
	// * blue
	//
	// * blue_violet
	//
	// * cadet_blue
	//
	// * chocolate
	//
	// * coral
	//
	// * dodger_blue
	//
	// * firebrick
	//
	// * golden_rod
	//
	// * green
	//
	// * hot_pink
	//
	// * orange_red
	//
	// * red
	//
	// * sea_green
	//
	// * spring_green
	//
	// * yellow_green
	//
	// Turbo and Prime users may specify a named color or a Hex color code like #9146FF.
	Color string
}

// UpdateUserChatColor implements https://dev.twitch.tv/docs/api/reference#update-user-chat-color
//
// Updates the color used for the user’s name in chat.
func (c *helixClient) UpdateUserChatColor(ctx context.Context, req *UpdateUserChatColorRequest) error {
	values := url.Values{}
	values.Set("user_id", req.UserID)
	values.Set("color", req.Color)

	return client.WithoutBody(c.Request(&client.RequestConfig{
		Method:  http.MethodPut,
		URL:     chatColorPath,
		Headers: c.headers(req.RequestOptions),
		Query:   values,
	}).Do(ctx))
}
