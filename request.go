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
	logFunc func(ctx context.Context, req *http.Request, reqBodyBytes, respBodyBytes []byte, requestTime, responseTime time.Time, statusCode int, respHeaders http.Header) error
}

func New(
	timeout time.Duration,
	logFunc func(ctx context.Context, req *http.Request, reqBodyBytes, respBodyBytes []byte, requestTime, responseTime time.Time, statusCode int, respHeaders http.Header) error,
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
	if data, ok := result.Data.([]interface{}); ok {
		msg = data[0].(string)
	} else {
		msg = result.Data.(string)
	}

	return errors.New(msg)
}
