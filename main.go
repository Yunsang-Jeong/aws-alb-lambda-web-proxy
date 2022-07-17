package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Yunsang-Jeong/aws-alb-lambda-web-proxy/app"
)

func main() {
	lambda.Start(app.Run)
}
