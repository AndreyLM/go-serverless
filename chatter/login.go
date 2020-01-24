package main

import (
	"context"

	"github.com/andreylm/go-serverless/chatter/chatsess"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-lambda-go/lambda"
)

// Event - event
type Event struct {
	Username string
	Password string
}

const jobName = "Login"

// Response - response
type Response struct {
	Job    string
	Sessid string
	Err    string
}

func newResponse(sessid, err string) Response {
	return Response{
		Job:    jobName,
		Sessid: sessid,
		Err:    err,
	}
}

func handler(c context.Context, ev Event) (Response, error) {
	sess := session.Must(session.NewSession())
	u, err := chatsess.GetDBUserPass(ev.Username, ev.Password, sess)
	if err != nil {
		return newResponse("", "GetDBUserPass: "+err.Error()), nil
	}

	lg := chatsess.NewLogin(u.Username)
	err = lg.Put(sess)
	if err != nil {
		return newResponse("", "Put: "+err.Error()), nil
	}

	return newResponse(lg.Sessid, ""), nil
}

func main() {
	lambda.Start(handler)
}
