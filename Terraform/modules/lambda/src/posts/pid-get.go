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

func pidGet(request events.APIGatewayProxyRequest) (any, int, error) {
	id := request.PathParameters["id"]
	if id == "" {
		return nil, 400, fmt.Errorf("id is required")
	}

	return getPost(id)
}

func getPost(id string) (any, int, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, 500, err
	}
	db := dynamodb.New(sess)

	tn := os.Getenv("POSTS_TABLE_NAME")
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tn),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	result, err := db.GetItem(input)
	if err != nil {
		return nil, 500, err
	}

	post := Post{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &post)
	if err != nil {
		return nil, 500, err
	}
	if post.Id == nil {
		return nil, 404, fmt.Errorf("post not found")
	}
	response := struct {
		Post Post `json:"post"`
	}{
		Post: post,
	}
	return response, 200, nil
}
