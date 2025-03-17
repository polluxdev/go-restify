package webapi

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type contextKey string

const requestIDKey contextKey = "requestId"

type loggingTransport struct {
	Transport http.RoundTripper
	LogFunc   LogFunc
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Transport == nil {
		t.Transport = http.DefaultTransport
	}

	ctx := req.Context()
	requestTime := time.Now()

	var reqBodyBytes []byte
	if req.Body != nil {
		reqBodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(reqBodyBytes))
	}

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewReader(responseBodyBytes))

	responseTime := time.Now()

	requestId, _ := ctx.Value(requestIDKey).(string)

	go t.LogFunc(ctx, requestId, req, reqBodyBytes, responseBodyBytes, requestTime, responseTime, resp.StatusCode, resp.Header)

	return resp, nil
}
