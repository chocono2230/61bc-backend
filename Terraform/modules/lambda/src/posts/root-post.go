package posts

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type RequestBody struct {
	UserId  *string `json:"userId"`
	ReplyId *string `json:"replyId"`
	Content *struct {
		Comment *string `json:"comment"`
	} `json:"content"`
}

func post(request events.APIGatewayProxyRequest) (any, int, error) {
	body := request.Body
	rb := RequestBody{}
	err := json.Unmarshal([]byte(body), &rb)
	if err != nil {
		return nil, 400, fmt.Errorf("request body json unmarshal error")
	}
	if rb.UserId == nil || *rb.UserId == "" {
		return nil, 400, fmt.Errorf("user id is required")
	}
	if rb.Content == nil {
		return nil, 400, fmt.Errorf("content is required")
	}

	switch {
	case rb.Content.Comment != nil && *rb.Content.Comment != "":
		return createPost(rb)
	default:
		return nil, 400, fmt.Errorf("content is required")
	}
}

func createPost(requestBody RequestBody) (any, int, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, 500, err
	}
	db := dynamodb.New(sess)

	id := uuid.New().String()
	timestamp := time.Now().Unix()
	var replyId *string
	if requestBody.ReplyId != nil && *requestBody.ReplyId != "" {
		replyId = requestBody.ReplyId
	}
	content := requestBody.Content
	var gsiSKey string = "SKEY"
	post := Post{
		Id:        &id,
		UserId:    requestBody.UserId,
		Timestamp: &timestamp,
		GsiSKey:   &gsiSKey,
		ReplyId:   replyId,
	}
	if content.Comment != nil && *content.Comment != "" {
		post.Content = &struct {
			Comment *string `dynamodbav:"comment" json:"comment"`
		}{
			Comment: content.Comment,
		}
	}

	inputAV, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		return nil, 500, err
	}
	tn := os.Getenv("POSTS_TABLE_NAME")
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tn),
		Item:      inputAV,
	}
	_, err = db.PutItem(input)
	if err != nil {
		return nil, 500, err
	}

	return post, 201, nil
}
