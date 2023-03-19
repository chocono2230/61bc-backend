package posts

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type UpdateReactionRequest = struct {
	Type    *string `json:"type"`
	Payload *any    `json:"payload"`
}

func pidPatch(request events.APIGatewayProxyRequest) (any, int, error) {
	body := request.Body
	ur := UpdateReactionRequest{}
	err := json.Unmarshal([]byte(body), &ur)
	if err != nil {
		return nil, 400, fmt.Errorf("request body json unmarshal error")
	}
	id := request.PathParameters["id"]
	if ur.Type == nil || id == "" {
		return nil, 400, fmt.Errorf("id, postId and type are required")
	}

	return updateReaction(id, *ur.Type, ur.Payload)
}

func updateReaction(id string, ttype string, payload *any) (any, int, error) {
	switch ttype {
	case "like":
		return updateLike(id)
	}
	return nil, 400, fmt.Errorf("invalid patch type")
}

func updateLike(id string) (any, int, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, 500, err
	}
	db := dynamodb.New(sess)

	basepost, st, err := internalGetPost(id)
	if err != nil {
		return nil, st, err
	}
	var update expression.UpdateBuilder
	if basepost.Reactions == nil {
		update = expression.UpdateBuilder{}.Set(expression.Name("reactions"), expression.Value(1))
	} else {
		update = expression.UpdateBuilder{}.Add(expression.Name("reactions"), expression.Value(1))
	}
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return nil, 500, err
	}

	tn := os.Getenv("POSTS_TABLE_NAME")
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tn),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              aws.String("ALL_NEW"),
	}

	result, err := db.UpdateItem(input)
	if err != nil {
		return nil, 500, err
	}
	post := Post{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &post)
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
	return &response, 200, nil
}
