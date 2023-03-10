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
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	DisplayName *string `json:"displayName"`
	Identity    *string `json:"identity"`
}

func post(request events.APIGatewayProxyRequest) (any, int, error) {
	body := request.Body
	cr := CreateUserRequest{}
	err := json.Unmarshal([]byte(body), &cr)
	if err != nil {
		return nil, 400, fmt.Errorf("request body json unmarshal error")
	}
	if cr.DisplayName == nil || cr.Identity == nil || *cr.DisplayName == "" || *cr.Identity == "" {
		return nil, 400, fmt.Errorf("displayName and identity are required")
	}

	return createUser(cr)
}

func createUser(cr CreateUserRequest) (any, int, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, 500, err
	}
	db := dynamodb.New(sess)
	tableName := os.Getenv("USERS_TABLE_NAME")

	id := uuid.New().String()
	user := User{
		Id:          &id,
		DisplayName: cr.DisplayName,
		Identity:    cr.Identity,
	}

	iav, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, 500, err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      iav,
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
	return rs, 201, nil
}
