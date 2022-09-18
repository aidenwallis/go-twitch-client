package helix

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aidenwallis/go-twitch-client/internal/testutils"
	"github.com/aidenwallis/go-twitch-client/internal/testutils/assert"
)

func testClient(it testutils.RoundTripInterceptor) Client {
	return NewClient(&ClientOptions{
		ClientID:  fakeClientID,
		Transport: testutils.Middleware(it),
	})
}

const (
	fakeClientID = "clientID"
	fakeToken    = "abc123"
)

func requestOptions() *RequestOptions {
	return &RequestOptions{
		Token: fakeToken,
	}
}

func assertToken(t *testing.T, req *http.Request) {
	assert.Equal(t, fakeClientID, req.Header.Get("Client-ID"), "client id must be present")
	assert.Equal(t, "Bearer "+fakeToken, req.Header.Get("authorization"), "token must be present")
}

func TestDefaultOptions(t *testing.T) {
	// a hacky test, but ensures we don't panic if we have a nil client options passed
	ctx := context.Background()
	c := NewClient(nil)
	c.(*helixClient).Transport = testutils.Middleware(func(req *http.Request) *testutils.Response {
		assert.Equal(t, "", req.Header.Get("Client-ID"))
		return testutils.EmptyResponse(http.StatusOK)
	})

	assert.NoError(t, c.BlockUser(ctx, &BlockUserRequest{TargetUserID: "id"}))
}

func TestAccessTokenLoader(t *testing.T) {
	ctx := context.Background()
	expected := errors.New("expected")

	t.Run("error", func(t *testing.T) {
		c := NewClient(&ClientOptions{
			AccessTokenLoader: func(ctx context.Context) (string, error) {
				return "", expected
			},
			Transport: testutils.Middleware(func(req *http.Request) *testutils.Response {
				assertToken(t, req)
				return testutils.EmptyResponse(http.StatusOK)
			}),
		})
		if c.BlockUser(ctx, &BlockUserRequest{TargetUserID: "id"}) != expected {
			t.Error("expected error to be thrown")
		}
	})

	t.Run("success", func(t *testing.T) {
		c := NewClient(&ClientOptions{
			AccessTokenLoader: func(ctx context.Context) (string, error) {
				return fakeToken, nil
			},
			ClientID: fakeClientID,
			Transport: testutils.Middleware(func(req *http.Request) *testutils.Response {
				assertToken(t, req)
				return testutils.EmptyResponse(http.StatusOK)
			}),
		})
		assert.NoError(t, c.BlockUser(ctx, &BlockUserRequest{TargetUserID: "id"}))
	})
}
