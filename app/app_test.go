package app_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Yunsang-Jeong/aws-alb-lambda-web-proxy/app"
	"github.com/aws/aws-lambda-go/events"
)

func Test_Run_default(t *testing.T) {
	default_method := "GET"
	default_path := "/"
	default_headers := map[string]string{
		"host":              "www.naver.com",
		"user-agent":        "golang-test",
		"x-amzn-trace-id":   "Root=1-01234567-012345678901234567890123",
		"x-forwarded-for":   "1.2.3.4",
		"x-forwarded-proto": "https",
	}

	resp, _ := app.Run(context.TODO(), events.ALBTargetGroupRequest{
		HTTPMethod:                      default_method,
		Path:                            default_path,
		QueryStringParameters:           map[string]string{},
		MultiValueQueryStringParameters: map[string][]string{},
		Headers:                         default_headers,
		MultiValueHeaders:               map[string][]string{},
		RequestContext:                  events.ALBTargetGroupRequestContext{},
		IsBase64Encoded:                 false,
		Body:                            "",
	})

	fmt.Println(resp)
}
