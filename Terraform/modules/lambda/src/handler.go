package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chocono2230/61bc-backend/lambda/healthcheck"
	"github.com/chocono2230/61bc-backend/lambda/posts"
)

func jsonResponse(body any, statusCode int, err error) (events.APIGatewayProxyResponse, error) {
	header := map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Headers": "*",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "*",
	}

	if err != nil {
		body = err.Error()
	}

	jsonBytes, err2 := json.Marshal(body)
	if err2 != nil {
		jsonBytes, err3 := json.Marshal("respons body json marshal error")
		if err3 != nil {
			return events.APIGatewayProxyResponse{
				Body:       "respons body json marshal error",
				StatusCode: 500,
				Headers:    header,
			}, err3
		}
		return events.APIGatewayProxyResponse{
			Body:       string(jsonBytes),
			StatusCode: 500,
			Headers:    header,
		}, err2
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: statusCode,
		Headers:    header,
	}, err
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := request.Path
	trm := strings.Trim(path, "/")
	pathArray := strings.Split(trm, "/")
	if len(pathArray) < 1 {
		return jsonResponse(nil, 400, fmt.Errorf("root path is not allowed"))
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
		err = fmt.Errorf("resource root is not allowed")
		statusCode = 400
	}
	return jsonResponse(res, statusCode, err)
}

func main() {
	lambda.Start(HandleRequest)
}
