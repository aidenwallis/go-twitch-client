package helix

import "github.com/aidenwallis/go-twitch-client/internal/testutils"

func testClient(it testutils.RoundTripInterceptor) Client {
	return NewClient(&ClientOptions{
		Transport: testutils.Middleware(it),
	})
}

func requestOptions() *RequestOptions {
	return &RequestOptions{
		Token: "abc123",
	}
}
