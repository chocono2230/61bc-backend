package users

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func get(request events.APIGatewayProxyRequest) (any, int, error) {
	return getAllUser()
}

func getAllUser() (any, int, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, 500, err
	}
	db := dynamodb.New(sess)

	tn := os.Getenv("USERS_TABLE_NAME")
	input := &dynamodb.ScanInput{
		TableName: aws.String(tn),
		ExpressionAttributeNames: map[string]*string{
			"#id":          aws.String("id"),
			"#displayName": aws.String("displayName"),
		},
		ProjectionExpression: aws.String("id, #displayName"),
	}
	result, err := db.Scan(input)
	if err != nil {
		return nil, 500, err
	}

	users := []PublicUser{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, 500, err
	}
	response := struct {
		Users []PublicUser `json:"users"`
	}{
		Users: users,
	}
	return response, 200, nil
}
