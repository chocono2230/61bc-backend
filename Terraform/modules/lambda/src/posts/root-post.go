package posts

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type RequestBody struct {
	UserId  *string `json:"userId"`
	ReplyId *string `json:"replyId"`
	Content *struct {
		Comment *string `json:"comment"`
	} `json:"content"`
}

func post(request events.APIGatewayProxyRequest) (string, int, error) {
	body := request.Body
	rb := RequestBody{}
	err := json.Unmarshal([]byte(body), &rb)
	if err != nil {
		return "request body json unmarshal error", 400, nil
	}
	if rb.UserId == nil || *rb.UserId == "" {
		return "user id is required", 400, nil
	}
	if rb.Content == nil {
		return "content is required", 400, nil
	}

	switch {
	case rb.Content.Comment != nil && *rb.Content.Comment != "":
		return fmt.Sprintf("comment: %s", *rb.Content.Comment), 200, nil
	default:
		return "content is required", 400, nil
	}
}
