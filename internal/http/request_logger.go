package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type RequestLogger struct{}

func (rl *RequestLogger) RoundTrip(r *http.Request) (*http.Response, error) {
	path, err := url.PathUnescape(r.URL.Path)
	if err != nil {
		path = r.URL.Path
	}

	rawQuery, err := url.QueryUnescape(r.URL.RawQuery)
	if err != nil {
		rawQuery = r.URL.RawQuery
	}

	if rawQuery != "" {
		path += "?" + rawQuery
	}

	fields := map[string]any{
		"request": fmt.Sprintf("%s %s", r.Method, path),
	}
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil {
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			fields["request_body"] = string(bodyBytes)
		}
	}

	startTime := time.Now()

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		fields["error"] = err.Error()
		tflog.Debug(r.Context(), "HTTP request failed", fields)

		return nil, err
	}

	duration := time.Since(startTime)
	fields["response_status"] = resp.Status
	fields["duration"] = duration
	tflog.Debug(r.Context(), "HTTP request completed", fields)

	return resp, nil
}
