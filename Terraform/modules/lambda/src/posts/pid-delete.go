package posts

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chocono2230/61bc-backend/lambda/users"
)

type DeletePostRequest = struct {
	UserId   *string `json:"userId"`
	Identity *string `json:"identity"`
}

func pidDelete(request events.APIGatewayProxyRequest) (any, int, error) {
	body := request.Body
	dr := DeletePostRequest{}
	err := json.Unmarshal([]byte(body), &dr)
	if err != nil {
		return nil, 400, fmt.Errorf("request body json unmarshal error")
	}
	id := request.PathParameters["id"]
	if dr.Identity == nil || dr.UserId == nil || id == "" {
		return nil, 400, fmt.Errorf("id, userId and identity are required")
	}

	return deletePost(id, dr)
}

func deletePost(id string, dr DeletePostRequest) (any, int, error) {
	ps, st, err := internalGetPost(id)
	if err != nil || ps == nil {
		return nil, st, err
	}
	if *ps.UserId != *dr.UserId {
		return nil, 400, fmt.Errorf("user id is not matched")
	}
	_, st, err = users.UserVerification(*dr.UserId, *dr.Identity)
	if err != nil {
		return nil, st, err
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, 500, err
	}
	db := dynamodb.New(sess)

	tn := os.Getenv("POSTS_TABLE_NAME")
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tn),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	_, err = db.DeleteItem(input)
	if err != nil {
		return nil, 500, err
	}

	return nil, 204, nil
}
