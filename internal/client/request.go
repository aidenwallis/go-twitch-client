package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type Request struct {
	requestBody    []byte
	err            error
	client         *Client
	headers        http.Header
	headersFactory func(context.Context) (http.Header, error)
	query          url.Values
	method         string
	url            string
}

type RequestConfig struct {
	Method string

	// URL is assumed to not have any query encoded values, and it does not check
	// in the Do function, this is a deliberate assumption as this only lives in an
	// internal package. Please make sure your URLs do not include query strings,
	// use Query instead.
	URL string

	// Headers are all HTTP headers to attach to the request.
	Headers func(context.Context) (http.Header, error)

	// Query are all query string keys/values to attach to the URL.
	Query url.Values
}

// Request creates a new request for the client.
func (c *Client) Request(conf *RequestConfig) *Request {
	return &Request{
		client:         c,
		headersFactory: conf.Headers,
		method:         conf.Method,
		query:          conf.Query,
		url:            conf.URL,
	}
}

// BodyJSON encodes a JSON body and attaches it to the request
func (r *Request) BodyJSON(body interface{}) *Request {
	if r.err != nil {
		return r
	}

	r.requestBody, r.err = json.Marshal(body)
	if r.err != nil {
		return r
	}

	r.headers.Set("Content-Type", "application/json; charset=utf-8")
	return r
}

func (r *Request) Do(ctx context.Context) *Response {
	url := r.url
	if len(r.query) > 0 {
		url += "?" + r.query.Encode()
	}

	h, err := r.headersFactory(ctx)
	if err != nil {
		return &Response{Response: nil, err: err}
	}

	req, err := http.NewRequestWithContext(ctx, r.method, url, bytes.NewReader(r.requestBody))
	if err != nil {
		return &Response{Response: nil, err: err}
	}

	for key, vs := range r.headers {
		for _, v := range vs {
			req.Header.Add(key, v)
		}
	}

	for key, vs := range h {
		for _, v := range vs {
			req.Header.Add(key, v)
		}
	}

	res, err := r.client.Do(req)
	return &Response{Response: res, err: err}
}
