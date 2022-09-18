package helix

import (
	"context"
	"net/http"
	"testing"

	"github.com/aidenwallis/go-twitch-client/internal/testutils"
	"github.com/aidenwallis/go-twitch-client/internal/testutils/assert"
)

func TestBlockUser(t *testing.T) {
	ctx := context.Background()
	in := &BlockUserRequest{
		RequestOptions: requestOptions(),

		TargetUserID:  "targetUserID",
		SourceContext: "chat",
		Reason:        "some reason",
	}

	c := testClient(func(req *http.Request) *testutils.Response {
		assert.Equal(t, in.TargetUserID, req.URL.Query().Get("target_user_id"))
		assert.Equal(t, in.SourceContext, req.URL.Query().Get("source_context"))
		assert.Equal(t, in.Reason, req.URL.Query().Get("reason"))
		assert.Equal(t, userBlocksPath, testutils.WithoutQuery(req.URL))
		return testutils.EmptyResponse(http.StatusNoContent)
	})

	assert.NoError(t, c.BlockUser(ctx, in))
}
