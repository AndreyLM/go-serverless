package main

import (
	"context"

	"github.com/andreylm/go-serverless/chatter/chatsess"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Event - event
type Event struct {
	Sessid   string
	LastID   string
	LastTime string
}

const jobName = "read"

// Response - response
type Response struct {
	Job   string
	Err   string
	Chats []chatsess.Chat
}

func newResponse(chats, err string) Response {
	return Response{
		Job:   jobName,
		Err:   err,
		Chats: chats,
	}
}

func handler(c context.Context, ev Event) (Response, error) {
	sess := session.Must(session.NewSession())
	
	_, err := chatsess.GetLogin(ev.Sessid, sess)
	if err != nil {
		return newResponse(nil, "Not login: "+err.Error()), nil
	}

}

func main() {
	lambda.Start(handler)
}
