package testutils

import (
	"net/url"
	"strings"
)

// WithoutQuery returns the URL without a query string appended to it
func WithoutQuery(u *url.URL) string {
	return strings.TrimSuffix(u.String(), "?"+u.Query().Encode())
}
