package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3Client *s3.S3

func handleRequest(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	log.Printf("Processing request data for traceId %s.\n", request.Headers["x-amzn-trace-id"])
	log.Printf("Body size = %d.\n", len(request.Body))

	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("    %s: %s\n", key, value)
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(&request)
	if err != nil {
		panic(err)
	}

	log.Printf("Storing response...")

	params := &s3.PutObjectInput{
		Bucket: aws.String("alb-responder"),
		Key:    aws.String(request.Path),
		Body:   bytes.NewReader(body.Bytes()),
	}

	output, err := s3Client.PutObject(params)
	if err != nil {
		return events.ALBTargetGroupResponse{}, err
	}

	log.Printf("Written request: %s\n", output)

	return events.ALBTargetGroupResponse{
		Body:              "OK body",
		StatusCode:        200,
		StatusDescription: "200 OK",
		IsBase64Encoded:   false,
		Headers:           map[string]string{},
	}, nil
}

func main() {
	sess := session.Must(session.NewSession())
	s3Client = s3.New(sess, &aws.Config{
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		Region:           aws.String("eu-west-1"),
		Endpoint:         aws.String("http://" + os.Getenv("LOCALSTACK_HOSTNAME") + ":" + os.Getenv("EDGE_PORT")),
		S3ForcePathStyle: aws.Bool(true),
	})

	lambda.Start(handleRequest)
}
