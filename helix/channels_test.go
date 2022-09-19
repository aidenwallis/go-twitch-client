package helix

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/aidenwallis/go-twitch-client"
	"github.com/aidenwallis/go-twitch-client/internal/testutils"
	"github.com/aidenwallis/go-twitch-client/internal/testutils/assert"
)

func TestGetChannelInformation(t *testing.T) {
	ctx := context.Background()
	in := &GetChannelInformationRequest{
		RequestOptions: requestOptions(),
		BroadcasterIDs: []string{"1", "2"},
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, channelsPath+"?broadcaster_id=1&broadcaster_id=2", req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &GetChannelInformationResponse{
			Data: []*Channel{
				{
					BroadcasterID:       "1",
					BroadcasterLogin:    "one",
					BroadcasterName:     "one",
					BroadcasterLanguage: "en",
					GameName:            "Just Chatting",
					GameID:              "123",
					Title:               "title",
					Delay:               0,
				},
				{
					BroadcasterID:       "2",
					BroadcasterLogin:    "two",
					BroadcasterName:     "two",
					BroadcasterLanguage: "en",
					GameName:            "Just Chatting",
					GameID:              "123",
					Title:               "title",
					Delay:               0,
				},
			},
		})
	})

	resp, err := c.GetChannelInformation(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Data))
}

func TestModifyChannelInformation(t *testing.T) {
	ctx := context.Background()
	in := &ModifyChannelInformationRequest{
		RequestOptions: requestOptions(),
		BroadcasterID:  "1",
		Title:          twitch.Pointer("title"),
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assert.Equal(t, http.MethodPatch, req.Method)
		assert.Equal(t, channelsPath+"?broadcaster_id=1", req.URL.String())
		assert.Equal(t, `{"title":"title"}`, testutils.DecodeRawBody(t, req))
		return testutils.EmptyResponse(http.StatusNoContent)
	})

	assert.NoError(t, c.ModifyChannelInformation(ctx, in))
}

func TestGetChannelEditors(t *testing.T) {
	ctx := context.Background()
	in := &GetChannelEditorsRequest{
		RequestOptions: requestOptions(),
		BroadcasterID:  "1",
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assertToken(t, req)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, channelEditorsPath+"?broadcaster_id=1", req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &GetChannelEditorsResponse{
			Data: []*ChannelEditor{
				{
					UserID:    "1",
					UserName:  "one",
					CreatedAt: time.Now().UTC(),
				},
				{
					UserID:    "2",
					UserName:  "two",
					CreatedAt: time.Now().UTC(),
				},
			},
		})
	})

	resp, err := c.GetChannelEditors(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Data))
}
