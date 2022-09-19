package helix

import (
	"context"
	"net/http"
	"testing"

	"github.com/aidenwallis/go-twitch-client/internal/testutils"
	"github.com/aidenwallis/go-twitch-client/internal/testutils/assert"
)

func TestStartCommercial(t *testing.T) {
	ctx := context.Background()
	in := &StartCommercialRequest{
		RequestOptions: requestOptions(),
		BroadcasterID:  "broadcasterID",
		Length:         60,
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, commercialPath, req.URL.String())
		return testutils.JSONResponse(t, http.StatusOK, &StartCommercialResponse{
			Data: []*Commercial{
				{
					Length:     60,
					Message:    "",
					RetryAfter: 480,
				},
			},
		})
	})

	resp, err := c.StartCommercial(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Data))
}
