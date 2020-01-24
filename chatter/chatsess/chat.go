package chatsess

import (
	"html"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Chat - chat
type Chat struct {
	DateID   string
	Time     time.Time
	Username string
	Text     string
}

// NewChat - creates new chat
func NewChat(username, text string) Chat {
	return Chat{
		DateID:   time.Now().Format(DATE_FMT),
		Time:     time.Now(),
		Username: username,
		Text:     html.EscapeString(text),
	}
}

// ChatFromItem - chat
func ChatFromItem(item map[string]*dynamodb.AttributeValue) Chat {
	dateav := item["DateID"]
	timeav := item["Tmstp"]
	unameav := item["Username"]
	txav := item["Text"]

	return Chat{
		DateID:   *dateav.S,
		Time:     DBtoTime(timeav.N),
		Username: *unameav.S,
		Text:     *txav.S,
	}
}

// Put - puts to db
func (c Chat) Put(sess *session.Session) error {
	dbc := dynamodb.New(sess)

	_, err := dbc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("ch_chats"),
		Item: map[string]*dynamodb.AttributeValue{
			"DateID":   {S: aws.String(c.DateID)},
			"Tmstp":    {N: TimetoDB(c.Time)},
			"Username": {S: aws.String(c.Username)},
			"Text":     {S: aws.String(c.Text)},
		},
	})

	return err
}

// GetChat - gets chat
func GetChat(sess *session.Session) ([]Chat, error) {
	dbc := dynamodb.New(sess)

	dbres, err := dbc.Query(&dynamodb.QueryInput{
		TableName:              aws.String("ch_chats"),
		KeyConditionExpression: aws.String("DateID = :a"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {S: aws.String(time.Now().Format(DATE_FMT))},
		},
	})
	if err != nil {
		return nil, err
	}
	res := []Chat{}
	for _, v := range dbres.Items {
		res = append(res, ChatFromItem(v))
	}

	return res, nil
}

// GetChatAfter - gets chat
func GetChatAfter(DateID string, t time.Time, sess *session.Session) ([]Chat, error) {
	dbc := dynamodb.New(sess)

	dbres, err := dbc.Query(&dynamodb.QueryInput{
		TableName:              aws.String("ch_chats"),
		KeyConditionExpression: aws.String("DateID = :a"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {S: aws.String(time.Now().Format(DATE_FMT))},
		},
		ExclusiveStartKey: map[string]*dynamodb.AttributeValue{
			"DateID": {S: aws.String(DateID)},
			"Tmstp":  {N: TimetoDB(t)},
		},
	})
	if err != nil {
		return nil, err
	}
	res := []Chat{}
	for _, v := range dbres.Items {
		res = append(res, ChatFromItem(v))
	}

	return res, nil
}
