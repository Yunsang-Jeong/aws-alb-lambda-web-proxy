package main

import (
	"context"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joeshaw/envdecode"
	"github.com/pkg/errors"
)

type LambdaConfig struct {
	ProxyUrl       string `env:"PROXY_URL"`
	RequestTimeOut int64  `env:"REQUEST_TIME_OUT,default=5"`
}

type Lambda struct {
	ProxyUrl   string
	HttpClient *http.Client
}

type proxyResponse struct {
	statusCode        int
	statusDescription string
	body              string
}

func logging(label string, data interface{}) {
	switch data := data.(type) {
	case map[string]string:
		for key, value := range data {
			log.Printf("[%s] %s: %s\n", label, key, value)
		}
	case map[string][]string:
		for key, value := range data {
			log.Printf("[%s] %s: %s\n", label, key, value)
		}
	case []string:
		for _, element := range data {
			log.Printf("[%s] %s\n", label, element)
		}
	case []byte:
		log.Printf("[%s] %s\n", label, string(data[:]))
	case string:
		log.Printf("[%s] %s\n", label, data)
	case int:
		log.Printf("[%s] %d\n", label, data)
	default:
		log.Printf("[%s] Error to print this data\n", label)
	}
}

func generateResponse(statusCode int, statusDescription string, body string) events.ALBTargetGroupResponse {
	return events.ALBTargetGroupResponse{
		StatusCode:        statusCode,
		StatusDescription: statusDescription,
		Headers:           map[string]string{},
		Body:              body,
		IsBase64Encoded:   false,
	}
}

func proxy(cli *http.Client, url string, payload string) (*proxyResponse, error) {
	resp, err := cli.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		return nil, errors.Wrap(err, "fail to send post")
	}

	headers := map[string]string{}
	for key, value := range resp.Header {
		headers[key] = strings.Join(value, "; ")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "fail to read body from reponse")
	}

	return &proxyResponse{
		statusCode:        resp.StatusCode,
		statusDescription: resp.Status,
		body:              string(body[:]),
	}, nil
}

func (l *Lambda) Handler(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	var payload string

	if request.IsBase64Encoded {
		decodedBody, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			logging("Error", err.Error())
			logging("Recv-Path", request.Path)
			logging("Recv-Header", request.Headers)
			logging("Recv-Body", request.Body)

			return generateResponse(400, "fail to decode body", ""), nil
		}

		payload = string(decodedBody[:])
	} else {
		payload = request.Body
	}

	resp, err := proxy(l.HttpClient, l.ProxyUrl+request.Path, payload)
	if err != nil {
		logging("Error", err.Error())
		logging("Recv-Path", request.Path)
		logging("Recv-Header", request.Headers)
		logging("Recv-Body", request.Body)
		logging("Recv-Body-Decoded", payload)

		return generateResponse(400, "fail to proxy", ""), nil
	}

	if len(resp.body) < 1024*1024 {
		return generateResponse(resp.statusCode, resp.statusDescription, resp.body), nil
	} else {
		logging("Error", "Response is bigger than 1MB")

		return generateResponse(400, "Response is bigger than 1MB", ""), nil
	}
}

func main() {
	config := LambdaConfig{}
	if err := envdecode.Decode(&config); err != nil {
		panic(err)
	}

	l := Lambda{
		ProxyUrl: config.ProxyUrl,
		HttpClient: &http.Client{
			Timeout: time.Duration(config.RequestTimeOut) * time.Second,
		},
	}

	lambda.Start(l.Handler)
}
