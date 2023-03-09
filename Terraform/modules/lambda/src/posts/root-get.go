package posts

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func get(request events.APIGatewayProxyRequest) (any, int, error) {
	uid := request.QueryStringParameters["userid"]

	switch {
	case uid != "":
		return nil, 400, fmt.Errorf("that feature has not been implemented")
	default:
		return getAllPost()
	}
}

func getAllPost() (any, int, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, 500, err
	}
	db := dynamodb.New(sess)

	tn := os.Getenv("POSTS_TABLE_NAME")
	in := os.Getenv("POSTS_TABLE_GSI_NAME_ALL")
	input := &dynamodb.QueryInput{
		IndexName: aws.String(in),
		TableName: aws.String(tn),
		ExpressionAttributeNames: map[string]*string{
			"#gsiSKey": aws.String("gsiSKey"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":gsiSKey": {
				S: aws.String("SKEY"),
			},
		},
		KeyConditionExpression: aws.String("#gsiSKey = :gsiSKey"),
		ScanIndexForward:       aws.Bool(false),
	}
	result, err := db.Query(input)
	if err != nil {
		return nil, 500, err
	}

	posts := []Post{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &posts)
	if err != nil {
		return nil, 500, err
	}
	response := struct {
		Posts []Post `json:"posts"`
	}{
		Posts: posts,
	}
	return response, 200, nil
}
