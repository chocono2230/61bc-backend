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
	fu, st, err := getUserFromIdentity(*cr.Identity)
	if st == 200 {
		rs := struct {
			User User `json:"user"`
		}{
			User: *fu,
		}
		return rs, 200, nil
	}
	if st != 404 {
		return nil, st, err
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, 500, err
	}
	db := dynamodb.New(sess)

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

	tn := os.Getenv("USERS_TABLE_NAME")
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tn),
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

func getUserFromIdentity(identity string) (*User, int, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, 500, err
	}
	db := dynamodb.New(sess)

	tn := os.Getenv("USERS_TABLE_NAME")
	in := os.Getenv("USERS_TABLE_GSI_NAME_IDENTITY")
	input := &dynamodb.QueryInput{
		TableName: aws.String(tn),
		IndexName: aws.String(in),
		ExpressionAttributeNames: map[string]*string{
			"#identity": aws.String("identity"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":identity": {
				S: aws.String(identity),
			},
		},
		KeyConditionExpression: aws.String("#identity = :identity"),
	}

	output, err := db.Query(input)
	if err != nil {
		return nil, 500, err
	}

	if len(output.Items) == 0 {
		return nil, 404, fmt.Errorf("user not found")
	}

	user := User{}
	err = dynamodbattribute.UnmarshalMap(output.Items[0], &user)
	if err != nil {
		return nil, 500, err
	}

	return &user, 200, nil
}
