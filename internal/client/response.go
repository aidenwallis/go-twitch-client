package client

import (
	"encoding/json"
	"net/http"

	"github.com/aidenwallis/go-twitch-client"
)

// TwitchError implements the Twitch API error structure
type TwitchError struct {
	Message string `json:"message"`
}

// Response wraps the HTTP response
type Response struct {
	*http.Response
	err error
}

// WithBody takes the response and unmarshals a body from it
func WithBody[Body any](r *Response) (*Body, error) {
	defer safeCleanupBody(r)

	if err := handleError(r); err != nil {
		return nil, err
	}

	var body Body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

func safeCleanupBody(r *Response) {
	if r.err == nil {
		r.Body.Close()
	}
}

// WithoutBody only inspects the HTTP status and returns without a body
func WithoutBody(r *Response) error {
	defer safeCleanupBody(r)
	return handleError(r)
}

func handleError(r *Response) error {
	if r.err != nil {
		return r.err
	}

	if r.StatusCode >= 200 && r.StatusCode < 300 {
		// 2xx response
		return nil
	}

	var body TwitchError
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return err
	}

	return twitch.NewError(body.Message, r.StatusCode)
}
