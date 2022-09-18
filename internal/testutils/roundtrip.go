package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// RoundTripInterceptor will intercept the HTTP client call and let you send your own response
type RoundTripInterceptor func(r *http.Request) *Response

// RoundTripper allows you to intercept HTTP client calls for testing
type roundTripper struct {
	interceptor RoundTripInterceptor
}

// RoundTrip intercepts the http transport roundtrip and lets you stub a response
func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := r.interceptor(req)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return &http.Response{
		StatusCode: resp.Status,
		Body:       io.NopCloser(bytes.NewReader(resp.Body)),
	}, nil
}

// Response wraps the roundtripper responses
type Response struct {
	Body    []byte
	Error   error
	Headers http.Header
	Status  int
}

// SetHeader adds a new header to the response
func (r *Response) SetHeader(key, value string) *Response {
	r.Headers.Set(key, value)
	return r
}

// EmptyResponse creates a new empty response
func EmptyResponse(status int) *Response {
	return &Response{Status: status}
}

// JSONResponse creates a new JSON response
func JSONResponse(t *testing.T, status int, jsonBody interface{}) *Response {
	bs, err := json.Marshal(jsonBody)
	if err != nil {
		t.Error("building JSON response", err)
	}

	h := http.Header{}
	h.Set("Content-Type", "application/json; charset=utf-8")

	return &Response{Body: bs, Headers: h, Status: status}
}

// ErrorResponse returns a new error response
func ErrorResponse(err error) *Response {
	return &Response{Error: err}
}

// Middleware creates a new interceptor middleware for a http transport
func Middleware(it RoundTripInterceptor) http.RoundTripper {
	return &roundTripper{interceptor: it}
}
