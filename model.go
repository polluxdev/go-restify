package restify

import (
	"context"
	"net/http"
	"time"
)

type LogFunc func(ctx context.Context, reqID string, req *http.Request, reqBodyBytes, respBodyBytes []byte, requestTime, responseTime time.Time, statusCode int, respHeaders http.Header) error

type WebApiResponseSuccess struct {
	Status  string      `bson:"status" json:"status"`
	Message string      `bson:"message" json:"message"`
	Data    interface{} `bson:"data" json:"data"`
}

type WebApiResponseFailed struct {
	WebApiResponseSuccess
	Error interface{} `bson:"error" json:"error"`
}
