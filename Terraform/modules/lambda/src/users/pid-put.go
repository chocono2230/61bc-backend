package users

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type UpdateUserRequest = User

func pidPut(request events.APIGatewayProxyRequest) (any, int, error) {
	body := request.Body
	ur := UpdateUserRequest{}
	err := json.Unmarshal([]byte(body), &ur)
	if err != nil {
		return nil, 400, fmt.Errorf("request body json unmarshal error")
	}
	if ur.Id == nil || ur.DisplayName == nil || ur.Identity == nil || *ur.DisplayName == "" || *ur.Id == "" || *ur.Identity == "" {
		return nil, 400, fmt.Errorf("id, displayName and identity are required")
	}

	return updateUser(ur)
}

func updateUser(ur UpdateUserRequest) (any, int, error) {
	fu, st, err := getUserFromId(*ur.Id)
	if err != nil {
		return nil, st, err
	}
	if *fu.Identity != *ur.Identity {
		return nil, 400, fmt.Errorf("identity cannot be changed")
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, 500, err
	}
	db := dynamodb.New(sess)

	user := User{
		Id:          ur.Id,
		DisplayName: ur.DisplayName,
		Identity:    ur.Identity,
	}

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, 500, err
	}
	tn := os.Getenv("USERS_TABLE_NAME")
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tn),
	}
	_, err = db.PutItem(input)
	if err != nil {
		return nil, 500, err
	}
	rs := struct {
		User User `json:"user"`
	}{
		User: user,
	}
	return rs, 200, nil
}
