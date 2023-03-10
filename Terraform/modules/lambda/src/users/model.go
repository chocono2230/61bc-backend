package users

type User struct {
	Id          *string `dynamodbav:"id" json:"id"`
	DisplayName *string `dynamodbav:"displayName" json:"displayName"`
	Identity    *string `dynamodbav:"identity" json:"identity"`
}
