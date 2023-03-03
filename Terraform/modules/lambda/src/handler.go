package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chocono2230/61bc-backend/lambda/greeting"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	greeting.SayHello()
	return fmt.Sprintf("Hello %s!", name.Name), nil
}

func main() {
	fmt.Println("mainのV3-1だよ")
	lambda.Start(HandleRequest)
}
