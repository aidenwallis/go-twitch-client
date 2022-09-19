package helix

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/aidenwallis/go-twitch-client/internal/testutils"
	"github.com/aidenwallis/go-twitch-client/internal/testutils/assert"
)

func TestGetExtensionAnalytics(t *testing.T) {
	ctx := context.Background()
	in := &GetExtensionAnalyticsRequest{
		After:       "after",
		ExtensionID: "extensionID",
		Type:        "overview_v2",
		First:       20,
		StartedAt:   time.Now().UTC(),
		EndedAt:     time.Now().Add(time.Hour * 12).UTC(),
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assert.Equal(t, in.After, req.URL.Query().Get("after"))
		assert.Equal(t, in.ExtensionID, req.URL.Query().Get("extension_id"))
		assert.Equal(t, in.Type, req.URL.Query().Get("type"))
		assert.Equal(t, strconv.Itoa(in.First), req.URL.Query().Get("first"))
		assert.Equal(t, in.StartedAt.Format(time.RFC3339), req.URL.Query().Get("started_at"))
		assert.Equal(t, in.EndedAt.Format(time.RFC3339), req.URL.Query().Get("ended_at"))

		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, extensionAnalyticsPath, testutils.WithoutQuery(req.URL))
		return testutils.JSONResponse(t, http.StatusOK, &GetExtensionAnalyticsResponse{
			Data: []*ExtensionAnalytic{
				{
					ExtensionID: in.ExtensionID,
					URL:         "https://twitch.tv",
					Type:        in.Type,
					DateRange: ExtensionAnalyticDateRange{
						StartedAt: in.StartedAt,
						EndedAt:   in.EndedAt,
					},
				},
			},
		})
	})

	resp, err := c.GetExtensionAnalytics(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Data))
}

func TestGetGameAnalytics(t *testing.T) {
	ctx := context.Background()
	in := &GetGameAnalyticsRequest{
		After:     "after",
		GameID:    "gameID",
		Type:      "overview_v2",
		First:     20,
		StartedAt: time.Now().UTC(),
		EndedAt:   time.Now().Add(time.Hour * 12).UTC(),
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assert.Equal(t, in.After, req.URL.Query().Get("after"))
		assert.Equal(t, in.GameID, req.URL.Query().Get("game_id"))
		assert.Equal(t, in.Type, req.URL.Query().Get("type"))
		assert.Equal(t, strconv.Itoa(in.First), req.URL.Query().Get("first"))
		assert.Equal(t, in.StartedAt.Format(time.RFC3339), req.URL.Query().Get("started_at"))
		assert.Equal(t, in.EndedAt.Format(time.RFC3339), req.URL.Query().Get("ended_at"))

		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, gameAnalyticsPath, testutils.WithoutQuery(req.URL))
		return testutils.JSONResponse(t, http.StatusOK, &GetGameAnalyticsResponse{
			Data: []*GameAnalytic{
				{
					GameID: in.GameID,
					URL:    "https://twitch.tv",
					Type:   in.Type,
					DateRange: GameAnalyticDateRange{
						StartedAt: in.StartedAt,
						EndedAt:   in.EndedAt,
					},
				},
			},
		})
	})

	resp, err := c.GetGameAnalytics(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Data))
}
