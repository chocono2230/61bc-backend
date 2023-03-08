package main

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chocono2230/61bc-backend/lambda/healthcheck"
	"github.com/chocono2230/61bc-backend/lambda/posts"
)

func jsonResponse(body any, statusCode int) (events.APIGatewayProxyResponse, error) {
	header := map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Headers": "*",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "*",
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		jsonBytes, err := json.Marshal("respons body json marshal error")
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       "respons body json marshal error",
				StatusCode: 500,
				Headers:    header,
			}, nil
		}
		return events.APIGatewayProxyResponse{
			Body:       string(jsonBytes),
			StatusCode: 500,
			Headers:    header,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: statusCode,
		Headers:    header,
	}, nil
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := request.Path
	trm := strings.Trim(path, "/")
	pathArray := strings.Split(trm, "/")
	if len(pathArray) < 1 {
		return jsonResponse("root path is not allowed", 400)
	}

	var res any
	statusCode := 500
	var err error
	switch pathArray[0] {
	case "healthcheck":
		res, statusCode, err = healthcheck.Root(request)
	case "posts":
		res, statusCode, err = posts.Root(request)
	default:
		res = "resource root is not allowed"
		statusCode = 400
	}
	if err != nil {
		res = err.Error()
	}
	return jsonResponse(res, statusCode)
}

func main() {
	lambda.Start(HandleRequest)
}
