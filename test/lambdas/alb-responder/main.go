package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3Client *s3.S3

func handleRequest(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	fmt.Printf("Processing request data for traceId %s.\n", request.Headers["x-amzn-trace-id"])
	fmt.Printf("Body size = %d.\n", len(request.Body))

	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("    %s: %s\n", key, value)
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(&request)
	if err != nil {
		panic(err)
	}

	params := &s3.PutObjectInput{
		Bucket: aws.String("alb-responder"),
		Key:    aws.String(request.Path),
		Body:   bytes.NewReader(body.Bytes()),
	}

	_, err = s3Client.PutObject(params)
	if err != nil {
		panic(err)
	}

	return events.ALBTargetGroupResponse{
		Body:              request.Body,
		StatusCode:        200,
		StatusDescription: "200 OK",
		IsBase64Encoded:   false,
		Headers:           map[string]string{},
	}, nil
}

func main() {
	sess := session.Must(session.NewSession())
	s3Client = s3.New(sess, &aws.Config{
		Region: aws.String("eu-west-1"),
	})

	lambda.Start(handleRequest)
}
