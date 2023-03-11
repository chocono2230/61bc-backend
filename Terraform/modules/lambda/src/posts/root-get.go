package posts

import (
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
		return getPostFromUid(uid)
	default:
		return getAllPost()
	}
}

func getPostFromUid(uid string) (any, int, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, 500, err
	}
	db := dynamodb.New(sess)

	tn := os.Getenv("POSTS_TABLE_NAME")
	in := os.Getenv("POSTS_TABLE_GSI_NAME_USR")
	input := &dynamodb.QueryInput{
		IndexName: aws.String(in),
		TableName: aws.String(tn),
		ExpressionAttributeNames: map[string]*string{
			"#userId": aws.String("userId"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userId": {
				S: aws.String(uid),
			},
		},
		KeyConditionExpression: aws.String("#userId = :userId"),
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
