package main_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	proxy "github.com/Yunsang-Jeong/aws-alb-lambda-web-proxy"
	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	l := proxy.Lambda{
		ProxyUrl: "https://push.mattermost.com",
		HttpClient: &http.Client{
			Timeout: time.Duration(5) * time.Second,
		},
	}

	req := events.ALBTargetGroupRequest{
		HTTPMethod:                      "POST",
		Path:                            "/",
		QueryStringParameters:           map[string]string{},
		MultiValueQueryStringParameters: map[string][]string{},
		Headers:                         map[string]string{},
		MultiValueHeaders:               map[string][]string{},
		RequestContext:                  events.ALBTargetGroupRequestContext{},
		IsBase64Encoded:                 false,
		Body:                            "",
	}

	resp, err := l.Handler(context.TODO(), req)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(resp)
}
