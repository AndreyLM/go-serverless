package main

import (
	"context"

	"github.com/andreylm/go-serverless/chatter/chatsess"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Event - event
type Event struct {
	Sessid string
	Text   string
}

const jobName = "say"

// Response - response
type Response struct {
	Job string
	Err string
}

func newResponse(job, err string) Response {
	return Response{
		Job: jobName + " " + job,
		Err: err,
	}
}

func handler(c context.Context, ev Event) (Response, error) {
	sess := session.Must(session.NewSession())
	lg, err := chatsess.GetLogin(ev.Sessid, sess)
	if err != nil {
		return newResponse(ev.Text, err.Error()), nil
	}

	ch := chatsess.NewChat(lg.Username, ev.Text)
	if err = ch.Put(sess); err != nil {
		return newResponse(ev.Text, err.Error()), nil
	}

	return newResponse(ev.Text, ""), nil
}

func main() {
	lambda.Start(handler)
}
