package users

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func getUserFromId(id string) (*User, int, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, 500, err
	}
	db := dynamodb.New(sess)

	tn := os.Getenv("USERS_TABLE_NAME")
	input := &dynamodb.QueryInput{
		TableName: aws.String(tn),
		ExpressionAttributeNames: map[string]*string{
			"#id": aws.String("id"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(id),
			},
		},
		KeyConditionExpression: aws.String("#id = :id"),
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

func UserVerification(id string, identity string) (any, int, error) {
	user, status, err := getUserFromId(id)
	if err != nil {
		return nil, status, err
	}

	if user.Identity == nil || *user.Identity != identity {
		return nil, 401, fmt.Errorf("user verification failed")
	}

	return "OK", 200, nil
}
