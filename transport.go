package gowebapi

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

type loggingTransport struct {
	Transport http.RoundTripper
	LogFunc   LogFunc
	RequestID string
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

	newCtx := context.WithoutCancel(ctx)

	go t.LogFunc(newCtx, t.RequestID, req, reqBodyBytes, responseBodyBytes, requestTime, responseTime, resp.StatusCode, resp.Header)

	return resp, nil
}
