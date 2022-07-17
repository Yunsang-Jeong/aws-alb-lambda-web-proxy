package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	errors "github.com/pkg/errors"
)

const (
	TARGET_HOST = "https://www.google.com/"
)

type reponse struct {
	statusCode        int
	headers           map[string]string
	base64EncodedBody string
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
	case string:
		log.Printf("[%s] %s\n", label, data)
	case int:
		log.Printf("[%s] %d\n", label, data)
	default:
		log.Printf("[%s] Error to print this data\n", label)
	}
}

func generateReturn(err interface{}, statusDescription string, resp *reponse) (events.ALBTargetGroupResponse, error) {
	if err != nil {
		logging("Error", err.(string))
	}

	if resp == nil {
		return events.ALBTargetGroupResponse{
			StatusCode:        200,
			StatusDescription: statusDescription,
			Headers:           map[string]string{},
			Body:              "",
			IsBase64Encoded:   false,
		}, nil
	} else {
		return events.ALBTargetGroupResponse{
			StatusCode:        200,
			StatusDescription: statusDescription,
			Headers:           map[string]string{},
			Body:              resp.base64EncodedBody,
			IsBase64Encoded:   true,
		}, nil
	}
}

func retrieve_web_page(host string, path string) (*reponse, error) {
	url := fmt.Sprintf("https://%s%s", host, path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("fail to retrieve web page from %s", host))
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
	base64EncodedBody := base64.StdEncoding.EncodeToString(body)

	return &reponse{resp.StatusCode, headers, base64EncodedBody}, nil
}

func Run(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	logging("ALB-Header", request.Headers)
	logging("ALB-Body", fmt.Sprintf("size is %d", len(request.Body)))

	resp, err := retrieve_web_page(request.Headers["host"], request.Path)
	if err != nil {
		return generateReturn(err.Error(), "Fail to retrieve", nil)
	}
	logging("Reponse-Header", resp.headers)
	logging("Reponse-Body", fmt.Sprintf("size is %d", len(resp.base64EncodedBody)))
	logging("Reponse-Status", resp.statusCode)

	if len(resp.base64EncodedBody) < 1024*1024 {
		return generateReturn(nil, "OK", resp)
	} else {
		return generateReturn(nil, "Payload is too big", nil)
	}
}
