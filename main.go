package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func main() {
	help := false
	flag.BoolVar(&help, "help", help, "display CLI help")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	log.Printf("Starting funpro server...")
	http.ListenAndServe(":8080", http.HandlerFunc(handle))
	log.Printf("Exiting...")
}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got new request...")

	awsConfig := aws.NewConfig().
		WithRegion("eu-west-1").
		WithEndpoint("http://localstack:4566").
		WithCredentials(credentials.NewCredentials(
			&credentials.EnvProvider{},
		))
	sess := session.Must(session.NewSession(awsConfig))
	client := lambda.New(sess, awsConfig)

	params := &lambda.InvokeInput{
		FunctionName: aws.String("alb-responder"),
		Payload: []byte(`{
	"requestContext": {
		"elb": {
			"targetGroupArn": "arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/lambda-279XGJDqGZ5rsrHC2Fjr/49e9d65c45c6791a"
		}
	},
	"httpMethod": "GET",
	"path": "/lambda",
	"queryStringParameters": {
		"query": "1234ABCD"
	},
	"headers": {
		"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"accept-encoding": "gzip",
		"accept-language": "en-US,en;q=0.9",
		"connection": "keep-alive",
		"host": "lambda-alb-123578498.us-east-2.elb.amazonaws.com",
		"upgrade-insecure-requests": "1",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36",
		"x-amzn-trace-id": "Root=1-5c536348-3d683b8b04734faae651f476",
		"x-forwarded-for": "72.12.164.125",
		"x-forwarded-port": "80",
		"x-forwarded-proto": "http",
		"x-imforwards": "20"
	},
	"body": "",
	"isBase64Encoded": false
}`),
	}
	result, err := client.Invoke(params)
	if err != nil {
		panic(fmt.Sprintf("Could not invoke lambda %s: %v", "alb-responder", err))
	}

	w.WriteHeader(int(*result.StatusCode))
}
