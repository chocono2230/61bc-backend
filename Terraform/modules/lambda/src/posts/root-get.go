package posts

import (
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func get(request events.APIGatewayProxyRequest) (any, int, error) {
	uid := request.QueryStringParameters["userid"]
	eskId := request.QueryStringParameters["eskId"]
	eskTss := request.QueryStringParameters["eskTs"]
	eskTs := 0
	if eskTss != "" {
		eskTs, _ = strconv.Atoi(eskTss)
	}

	switch {
	case uid != "":
		return getPostFromUid(uid)
	default:
		return getAllPost(eskId, eskTs)
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

func getAllPost(eskId string, eskTs int) (any, int, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, 500, err
	}
	db := dynamodb.New(sess)

	tn := os.Getenv("POSTS_TABLE_NAME")
	in := os.Getenv("POSTS_TABLE_GSI_NAME_ALL")
	var esk map[string]*dynamodb.AttributeValue
	if eskId != "" {
		esk = map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(eskId),
			},
			"timestamp": {
				N: aws.String(strconv.Itoa(eskTs)),
			},
			"gsiSKey": {
				S: aws.String("SKEY"),
			},
		}
	}
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
		ExclusiveStartKey:      esk,
		Limit:                  aws.Int64(10),
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

	esk = result.LastEvaluatedKey
	if esk != nil {
		ts, _ := strconv.Atoi(*esk["timestamp"].N)
		response := struct {
			Posts []Post `json:"posts"`
			EskId string `json:"eskId"`
			EskTs int    `json:"eskTs"`
		}{
			Posts: posts,
			EskId: *esk["id"].S,
			EskTs: ts,
		}
		return response, 200, nil
	}

	response := struct {
		Posts []Post `json:"posts"`
	}{
		Posts: posts,
	}
	return response, 200, nil
}
