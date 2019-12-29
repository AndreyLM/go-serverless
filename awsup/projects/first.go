package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

// Event - event
type Event struct {
	Question string
}

// Response - lambda response
type Response struct {
	Question string
	Answer   string
}

func handler(e Event) (Response, error) {
	return Response{
		Question: e.Question,
		Answer:   "I don't know (third version)",
	}, nil
}

func main() {
	lambda.Start(handler)
}
