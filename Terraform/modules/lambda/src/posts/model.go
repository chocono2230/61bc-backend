package posts

type Post struct {
	Id          *string `dynamodbav:"id" json:"id"`
	UserId      *string `dynamodbav:"userId" json:"userId"`
	Timestamp   *int64  `dynamodbav:"timestamp" json:"timestamp"`
	GsiSKey     *string `dynamodbav:"gsiSKey" json:"gsiSKey"`
	ReplyId     *string `dynamodbav:"replyId" json:"replyId"`
	LastReplyId *string `dynamodbav:"lastReplyId" json:"lastReplyId"`
	Content     *struct {
		Comment *string `dynamodbav:"comment" json:"comment"`
		Image   *Image  `dynamodbav:"image" json:"image"`
	} `dynamodbav:"content" json:"content"`
	Reactions *[]any `dynamodbav:"reactions" json:"reactions"`
}

type Image struct {
	OriginId     *string `dynamodbav:"originId" json:"originId"`
	CompressedId *string `dynamodbav:"compressedId" json:"compressedId"`
}
