package logging

import (
	"math"
	"net/http"
	"time"
)

// RequestFields are set on an Entry by RequestLogs.
// They represent the standard set of fields to log when your service responds to an HTTP request.
type RequestFields struct {
	// These fields are mandatory
	Method     string `json:"method"`
	Url        string `json:"url"`
	StatusCode int    `json:"status"`
	DurationMs int    `json:"durationms"`
	// These ones are not mandatory
	BodySize     int64 `json:"bodysize,omitempty"`
	ResponseSize int64 `json:"responsesize,omitempty"`
	// The following fields are extensions to the spec, and are not mandatory
	RequestHeaders  http.Header `json:"requestheaders,omitempty"`
	RequestBody     interface{} `json:"requestbody,omitempty"`
	ResponseHeaders http.Header `json:"responseheaders,omitempty"`
	ResponseBody    interface{} `json:"responsebody,omitempty"`
}

func newRequestFields(l *logBase, r *http.Request, status int, responseSize int64, duration time.Duration) *RequestFields {
	return &RequestFields{
		Method:       r.Method,
		Url:          r.RequestURI,
		StatusCode:   status,
		BodySize:     r.ContentLength,
		ResponseSize: responseSize,
		DurationMs:   int(math.Ceil(float64(duration.Nanoseconds()) / 1000.0)),
	}
}
