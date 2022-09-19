package helix

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aidenwallis/go-twitch-client/internal/client"
)

// Analytics represents the Analytics namespace
type Analytics interface {
	// GetExtensionAnalytics implements https://dev.twitch.tv/docs/api/reference#get-extension-analytics
	//
	// Gets a URL that Extension developers can use to download analytics reports (CSV files) for their Extensions. The URL is valid for 5 minutes.
	//
	// If you specify a future date, the response will be “Report Not Found For Date Range.” If you leave both started_at and ended_at blank, the
	// API returns the most recent date of data.
	GetExtensionAnalytics(context.Context, *GetExtensionAnalyticsRequest) (*GetExtensionAnalyticsResponse, error)

	// GetGameAnalytics implements https://dev.twitch.tv/docs/api/reference#get-game-analytics
	//
	// Gets a URL that game developers can use to download analytics reports (CSV files) for their games. The URL is valid for 5 minutes. For
	// detail about analytics and the fields returned, see the Insights & Analytics guide: https://dev.twitch.tv/docs/insights
	//
	// The response has a JSON payload with a data field containing an array of games information elements and can contain a pagination field
	// containing information required to query for more streams.
	//
	// If you specify a future date, the response will be “Report Not Found For Date Range.” If you leave both started_at and ended_at blank, the
	// API returns the most recent date of data.
	GetGameAnalytics(context.Context, *GetGameAnalyticsRequest) (*GetGameAnalyticsResponse, error)
}

const extensionAnalyticsPath = "https://api.twitch.tv/helix/analytics/extensions"

// GetExtensionAnalyticsRequest represents the options passed to GetExtensionAnalytics
type GetExtensionAnalyticsRequest struct {
	*RequestOptions

	// After (optional) is the Cursor for forward pagination: tells the server where to start fetching the next set of results, in a
	// multi-page response. This applies only to queries without extension_id.
	//
	// If an extension_id is specified, it supersedes any cursor/offset combinations. The cursor value specified here is from the
	// pagination response field of a prior query.
	After string

	// ExtensionID (optional) is the Client ID value assigned to the extension when it is created. If this is specified, the returned URL
	// points to an analytics report for just the specified extension.
	//
	// If this is not specified, the response includes multiple URLs (paginated), pointing to separate analytics reports for each of the
	// authenticated user’s Extensions.
	ExtensionID string

	// Type (optional) of analytics report that is returned. Currently, this field has no affect on the response as there is only one report
	// type.
	//
	// If additional types were added, using this field would return only the URL for the specified report. Valid values: "overview_v2".
	Type string `json:"type"`

	// First (optional) is the Maximum number of objects to return. Maximum: 100. Default: 20.
	First int

	// StartedAt (optional) is the Starting date/time for returned reports, in RFC3339 format with the hours, minutes, and seconds zeroed out
	// and the UTC timezone: YYYY-MM-DDT00:00:00Z. This must be on or after January 31, 2018.
	//
	// If this is provided, ended_at also must be specified. If started_at is earlier than the default start date, the default date is used.
	// The file contains one row of data per day.
	StartedAt time.Time

	// EndedAt (optional) is the Ending date/time for returned reports, in RFC3339 format with the hours, minutes, and seconds zeroed out and
	// the UTC timezone: YYYY-MM-DDT00:00:00Z. The report covers the entire ending date; e.g., if 2018-05-01T00:00:00Z is specified, the
	// report covers up to 2018-05-01T23:59:59Z.
	//
	// If this is provided, started_at also must be specified. If ended_at is later than the default end date, the default date is used.
	// Default: 1-2 days before the request was issued (depending on report availability).
	EndedAt time.Time
}

// GetExtensionAnalyticsResponse represents the API response returned by GetExtensionAnalytics
type GetExtensionAnalyticsResponse struct {
	// Data represents a slice of ExtensionAnalytic
	Data []*ExtensionAnalytic `json:"data"`

	// Pagination represents Helix pagination data
	Pagination Pagination `json:"pagination"`
}

// ExtensionAnalytic represents a single extension analytic in Helix
type ExtensionAnalytic struct {
	// ExtensionID is the ID of the extension whose analytics data is being provided.
	ExtensionID string `json:"extension_id"`

	// URL to the downloadable CSV file containing analytics data. Valid for 5 minutes.
	URL string `json:"URL"`

	// Type of report.
	Type string `json:"type"`

	// DateRange are the date ranges for a given ExtensionAnalytic
	DateRange ExtensionAnalyticDateRange `json:"date_range"`
}

// ExtensionAnalyticDateRange represents the date range for a given ExtensionAnalytic
type ExtensionAnalyticDateRange struct {
	// StartedAt is the Report start date/time. Note this may differ from (be later than) the started_at value in the request;
	// the response value is the date when data for the extension is available.
	StartedAt time.Time `json:"started_at"`

	// EndedAt is the Report end date/time.
	EndedAt time.Time `json:"ended_at"`
}

// GetExtensionAnalytics implements https://dev.twitch.tv/docs/api/reference#get-extension-analytics
//
// Gets a URL that Extension developers can use to download analytics reports (CSV files) for their Extensions. The URL is valid for 5 minutes.
//
// If you specify a future date, the response will be “Report Not Found For Date Range.” If you leave both started_at and ended_at blank, the
// API returns the most recent date of data.
func (c *helixClient) GetExtensionAnalytics(ctx context.Context, req *GetExtensionAnalyticsRequest) (*GetExtensionAnalyticsResponse, error) {
	values := url.Values{}
	if req.After != "" {
		values.Set("after", req.After)
	}
	if req.ExtensionID != "" {
		values.Set("extension_id", req.ExtensionID)
	}
	if req.Type != "" {
		values.Set("type", req.Type)
	}
	if req.First > 0 {
		values.Set("first", strconv.Itoa(req.First))
	}
	if !req.StartedAt.IsZero() {
		values.Set("started_at", req.StartedAt.UTC().Format(time.RFC3339))
	}
	if !req.EndedAt.IsZero() {
		values.Set("ended_at", req.EndedAt.UTC().Format(time.RFC3339))
	}

	return client.WithBody[GetExtensionAnalyticsResponse](c.Request(&client.RequestConfig{
		Method:  http.MethodGet,
		URL:     extensionAnalyticsPath,
		Query:   values,
		Headers: c.headers(req.RequestOptions),
	}).Do(ctx))
}

const gameAnalyticsPath = "https://api.twitch.tv/helix/analytics/games"

// GetGameAnalyticsRequest represents the options passed to GetGameAnalytics
type GetGameAnalyticsRequest struct {
	*RequestOptions

	// After (optional) is the Cursor for forward pagination: tells the server where to start fetching the next set of results, in a multi-page
	// response. This applies only to queries without game_id.
	//
	// If a game_id is specified, it supersedes any cursor/offset combinations. The cursor value specified here is from the pagination response
	// field of a prior query.
	After string

	// GameID (optional). If this is specified, the returned URL points to an analytics report for just the specified game.
	//
	// If this is not specified, the response includes multiple URLs (paginated), pointing to separate analytics reports for each of the
	// authenticated user’s games.
	GameID string

	// Type (optional) of analytics report that is returned. Currently, this field has no affect on the response as there is only one report
	// type.
	//
	// If additional types were added, using this field would return only the URL for the specified report. Valid values: "overview_v2".
	Type string `json:"type"`

	// First (optional) is the Maximum number of objects to return. Maximum: 100. Default: 20.
	First int

	// StartedAt (optional) is the Starting date/time for returned reports, in RFC3339 format with the hours, minutes, and seconds zeroed out and
	// the UTC timezone: YYYY-MM-DDT00:00:00Z.
	//
	// If this is provided, ended_at also must be specified. If started_at is earlier than the default start date, the default date is used.
	// Default: 365 days before the report was issued. The file contains one row of data per day.
	StartedAt time.Time

	// EndedAt (optional) is the Ending date/time for returned reports, in RFC3339 format with the hours, minutes, and seconds zeroed out and the
	// UTC timezone: YYYY-MM-DDT00:00:00Z. The report covers the entire ending date; e.g., if 2018-05-01T00:00:00Z is specified, the report covers
	// up to 2018-05-01T23:59:59Z.
	//
	// If this is provided, started_at also must be specified. If ended_at is later than the default end date, the default date is used. Default:
	// 1-2 days before the request was issued (depending on report availability).
	EndedAt time.Time
}

// GetGameAnalyticsResponse defines the API response returned by GetGameAnalytics
type GetGameAnalyticsResponse struct {
	// Data represents a slice of GameAnalytic
	Data []*GameAnalytic `json:"data"`

	// Pagination represents Helix pagination data
	Pagination Pagination `json:"pagination"`
}

// GameAnalytic represents a game analytic entity in Helix
type GameAnalytic struct {
	// GameID is the ID of the game whose analytics data is being provided.
	GameID string `json:"game_id"`

	// URL to the downloadable CSV file containing analytics data. Valid for 5 minutes.
	URL string `json:"URL"`

	// Type of report.
	Type string `json:"type"`

	// DateRange are the date ranges for a given GameAnalytic
	DateRange GameAnalyticDateRange `json:"date_range"`
}

// GameAnalyticDateRange represents the date range for a given GameAnalytic
type GameAnalyticDateRange struct {
	// StartedAt is the Report start date/time.
	StartedAt time.Time `json:"started_at"`

	// EndedAt is the Report end date/time.
	EndedAt time.Time `json:"ended_at"`
}

// GetGameAnalytics implements https://dev.twitch.tv/docs/api/reference#get-game-analytics
//
// Gets a URL that game developers can use to download analytics reports (CSV files) for their games. The URL is valid for 5 minutes. For
// detail about analytics and the fields returned, see the Insights & Analytics guide: https://dev.twitch.tv/docs/insights
//
// The response has a JSON payload with a data field containing an array of games information elements and can contain a pagination field
// containing information required to query for more streams.
//
// If you specify a future date, the response will be “Report Not Found For Date Range.” If you leave both started_at and ended_at blank, the
// API returns the most recent date of data.
func (c *helixClient) GetGameAnalytics(ctx context.Context, req *GetGameAnalyticsRequest) (*GetGameAnalyticsResponse, error) {
	values := url.Values{}
	if req.After != "" {
		values.Set("after", req.After)
	}
	if req.GameID != "" {
		values.Set("game_id", req.GameID)
	}
	if req.Type != "" {
		values.Set("type", req.Type)
	}
	if req.First > 0 {
		values.Set("first", strconv.Itoa(req.First))
	}
	if !req.StartedAt.IsZero() {
		values.Set("started_at", req.StartedAt.UTC().Format(time.RFC3339))
	}
	if !req.EndedAt.IsZero() {
		values.Set("ended_at", req.EndedAt.UTC().Format(time.RFC3339))
	}

	return client.WithBody[GetGameAnalyticsResponse](c.Request(&client.RequestConfig{
		Method:  http.MethodGet,
		URL:     gameAnalyticsPath,
		Query:   values,
		Headers: c.headers(req.RequestOptions),
	}).Do(ctx))
}
