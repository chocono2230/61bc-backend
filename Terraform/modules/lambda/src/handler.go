package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	RequestMethod  string `json:"RequestMethod"`
	RequestBody    string `json:"RequestBody"`
	PathParameter  string `json:"PathParameter"`
	QueryParameter string `json:"QueryParameter"`
	Name           string `json:"Name"`
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// httpリクエストの情報を取得
	method := request.HTTPMethod
	body := request.Body
	pathParam := request.PathParameters["pathparam"]
	queryParam := request.QueryStringParameters["queryparam"]
	name := request.QueryStringParameters["name"]

	// レスポンスとして返すjson文字列を作る
	res := Response{
		RequestMethod:  method,
		RequestBody:    body,
		PathParameter:  pathParam,
		QueryParameter: queryParam,
		Name:           name,
	}
	jsonBytes, err := json.Marshal(res)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "respons json marshal error",
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
