package testutils

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/aidenwallis/go-twitch-client/internal/testutils/assert"
)

// DecodeJSONBody implements a testing util to let you easily decode request bodies in test middleware
func DecodeJSONBody[T any](t *testing.T, req *http.Request) *T {
	var out T
	assert.NoError(t, json.NewDecoder(req.Body).Decode(&out))
	return &out
}

// DecodeJSONBody implements a testing util to let you easily decode raw request bodies in test middleware
func DecodeRawBody(t *testing.T, req *http.Request) string {
	bs, err := io.ReadAll(req.Body)
	assert.NoError(t, err)
	return string(bs)
}
