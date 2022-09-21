package helix

import (
	"context"
	"net/http"
	"testing"

	"github.com/aidenwallis/go-twitch-client"
	"github.com/aidenwallis/go-twitch-client/internal/testutils"
	"github.com/aidenwallis/go-twitch-client/internal/testutils/assert"
)

func TestGetChannelEmotes(t *testing.T) {
	ctx := context.Background()
	in := &GetChannelEmotesRequest{
		RequestOptions: requestOptions(),
		BroadcasterID:  "1",
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, chatEmotesPath+"?broadcaster_id=1", req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &GetChannelEmotesResponse{
			Template: "https://twitch.tv/{{id}}_{{format}}_{{theme_mode}}_{{scale}}.png",
			Data: []*ChannelEmote{
				{
					ID:   "1",
					Name: "Kappa",
					Images: ChatEmoteImages{
						URL1x: "https://twitch.tv/1.png",
						URL2x: "https://twitch.tv/2.png",
						URL3x: "https://twitch.tv/3.png",
					},
					Tier:       "1000",
					EmoteType:  "subscriptions",
					EmoteSetID: "123",
					Format:     []string{"static"},
					Scale:      []string{"1.0", "2.0", "3.0"},
					ThemeMode:  []string{"light", "dark"},
				},
			},
		})
	})

	resp, err := c.GetChannelEmotes(ctx, in)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "https://twitch.tv/1_static_light_1.0.png", resp.EmoteURL("1", "static", "light", "1.0"))
}

func TestGetGlobalEmotes(t *testing.T) {
	ctx := context.Background()
	in := &GetGlobalEmotesRequest{
		RequestOptions: requestOptions(),
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, chatGlobalEmotesPath, req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &GetGlobalEmotesResponse{
			Template: "https://twitch.tv/{{id}}_{{format}}_{{theme_mode}}_{{scale}}.png",
			Data: []*GlobalEmote{
				{
					ID:   "1",
					Name: "Kappa",
					Images: ChatEmoteImages{
						URL1x: "https://twitch.tv/1.png",
						URL2x: "https://twitch.tv/2.png",
						URL3x: "https://twitch.tv/3.png",
					},
					Format:    []string{"static"},
					Scale:     []string{"1.0", "2.0", "3.0"},
					ThemeMode: []string{"light", "dark"},
				},
			},
		})
	})

	resp, err := c.GetGlobalEmotes(ctx, in)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "https://twitch.tv/1_static_light_1.0.png", resp.EmoteURL("1", "static", "light", "1.0"))
}

func TestGetEmoteSets(t *testing.T) {
	ctx := context.Background()
	in := &GetEmoteSetsRequest{
		RequestOptions: requestOptions(),
		EmoteSetIDs:    []string{"1", "2"},
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, chatEmoteSetsPath+"?emote_set_id=1&emote_set_id=2", req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &GetEmoteSetsResponse{
			Template: "https://twitch.tv/{{id}}_{{format}}_{{theme_mode}}_{{scale}}.png",
			Data: []*SetEmote{
				{
					ID:   "1",
					Name: "Kappa",
					Images: ChatEmoteImages{
						URL1x: "https://twitch.tv/1.png",
						URL2x: "https://twitch.tv/2.png",
						URL3x: "https://twitch.tv/3.png",
					},
					OwnerID:    "123",
					EmoteType:  "subscriptions",
					EmoteSetID: "123",
					Format:     []string{"static"},
					Scale:      []string{"1.0", "2.0", "3.0"},
					ThemeMode:  []string{"light", "dark"},
				},
			},
		})
	})

	resp, err := c.GetEmoteSets(ctx, in)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "https://twitch.tv/1_static_light_1.0.png", resp.EmoteURL("1", "static", "light", "1.0"))
}

func TestGetChannelChatBadges(t *testing.T) {
	ctx := context.Background()
	in := &GetChannelChatBadgesRequest{
		RequestOptions: requestOptions(),
		BroadcasterID:  "1",
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, channelChatBadges+"?broadcaster_id=1", req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &GetChannelChatBadgesResponse{
			Data: []*ChatBadge{
				{
					SetID: "moderator",
					Versions: []*ChatBadgeVersion{
						{
							ID:         "1",
							ImageURL1x: "https://twitch.tv/1.png",
							ImageURL2x: "https://twitch.tv/2.png",
							ImageURL4x: "https://twitch.tv/4.png",
						},
					},
				},
			},
		})
	})

	resp, err := c.GetChannelChatBadges(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Data))
}

func TestGetGlobalChatBadges(t *testing.T) {
	ctx := context.Background()
	in := &GetGlobalChatBadgesRequest{
		RequestOptions: requestOptions(),
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, globalChatBadges, req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &GetGlobalChatBadgesResponse{
			Data: []*ChatBadge{
				{
					SetID: "moderator",
					Versions: []*ChatBadgeVersion{
						{
							ID:         "1",
							ImageURL1x: "https://twitch.tv/1.png",
							ImageURL2x: "https://twitch.tv/2.png",
							ImageURL4x: "https://twitch.tv/4.png",
						},
					},
				},
			},
		})
	})

	resp, err := c.GetGlobalChatBadges(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Data))
}

func TestGetChatSettings(t *testing.T) {
	ctx := context.Background()
	in := &GetChatSettingsRequest{
		RequestOptions: requestOptions(),

		BroadcasterID: "1",
		ModeratorID:   "2",
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, chatSettingsPath+"?broadcaster_id=1&moderator_id=2", req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &GetChatSettingsResponse{
			Data: []*ChatSettings{
				{EmoteMode: true},
			},
		})
	})

	resp, err := c.GetChatSettings(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Data))
}

func TestUpdateChatSettings(t *testing.T) {
	ctx := context.Background()
	in := &UpdateChatSettingsRequest{
		RequestOptions: requestOptions(),

		BroadcasterID: "1",
		ModeratorID:   "2",

		EmoteMode: twitch.Pointer(true),
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodPatch, req.Method)
		assert.Equal(t, chatSettingsPath+"?broadcaster_id=1&moderator_id=2", req.URL.String())
		assert.Equal(t, `{"emote_mode":true}`, testutils.DecodeRawBody(t, req))
		return testutils.JSONResponse(t, http.StatusOK, &GetChatSettingsResponse{
			Data: []*ChatSettings{
				{EmoteMode: true},
			},
		})
	})

	resp, err := c.UpdateChatSettings(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Data))
}

func TestSendChatAnnouncement(t *testing.T) {
	ctx := context.Background()
	in := &SendChatAnnouncementRequest{
		RequestOptions: requestOptions(),

		BroadcasterID: "1",
		ModeratorID:   "2",
		Message:       "message",
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, chatAnnouncementsPath+"?broadcaster_id=1&moderator_id=2", req.URL.String())
		assert.Equal(t, `{"message":"message"}`, testutils.DecodeRawBody(t, req))
		return testutils.EmptyResponse(http.StatusNoContent)
	})

	assert.NoError(t, c.SendChatAnnouncement(ctx, in))
}

func TestGetUserChatColors(t *testing.T) {
	ctx := context.Background()
	in := &GetUserChatColorsRequest{
		RequestOptions: requestOptions(),

		UserIDs: []string{"1", "2"},
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, chatColorPath+"?user_id=1&user_id=2", req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &GetUserChatColorsResponse{
			Data: []*UserChatColor{
				{
					UserID:    "1",
					UserLogin: "one",
					UserName:  "one",
					Color:     "",
				},
				{
					UserID:    "2",
					UserLogin: "two",
					UserName:  "two",
					Color:     "",
				},
			},
		})
	})

	resp, err := c.GetUserChatColors(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Data))
}

func TestUpdateUserChatColor(t *testing.T) {
	ctx := context.Background()
	in := &UpdateUserChatColorRequest{
		RequestOptions: requestOptions(),

		UserID: "1",
		Color:  "purple",
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodPut, req.Method)
		assert.Equal(t, chatColorPath+"?color=purple&user_id=1", req.URL.String())
		return testutils.EmptyResponse(http.StatusNoContent)
	})

	assert.NoError(t, c.UpdateUserChatColor(ctx, in))
}
