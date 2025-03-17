package webapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type WebApiRequest interface {
	MakeRequest(ctx context.Context, method string, url string, headers map[string]string, body interface{}) (*http.Request, *http.Response, error)
	ParseErrorResponse(ctx context.Context, request *http.Request, response *http.Response) error
}

type webAPIRequest struct {
	timeout time.Duration
	logFunc LogFunc
}

func New(
	timeout time.Duration,
	logFunc LogFunc,
) WebApiRequest {
	return &webAPIRequest{
		timeout: timeout,
		logFunc: logFunc,
	}
}

func (w *webAPIRequest) MakeRequest(ctx context.Context, method string, url string, headers map[string]string, body interface{}) (*http.Request, *http.Response, error) {
	newCtx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, nil, err
	}

	request, err := http.NewRequestWithContext(newCtx, method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Add(k, v)
	}

	client := &http.Client{
		Timeout: w.timeout,
		Transport: &loggingTransport{
			Transport: http.DefaultTransport,
			LogFunc:   w.logFunc,
			RequestID: ctx.Value("requestId").(string),
		},
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}

	return request, response, nil
}

func (w *webAPIRequest) ParseErrorResponse(ctx context.Context, request *http.Request, response *http.Response) error {
	var result WebApiResponseFailed
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}

	var msg string
	if err, ok := result.Error.([]interface{}); ok {
		msg = err[0].(string)
	} else {
		msg = result.Error.(string)
	}

	return errors.New(msg)
}
