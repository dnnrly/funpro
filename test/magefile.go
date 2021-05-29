// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

func BuildAlbResponder() error {
	fmt.Println("Building alb-responder...")
	return sh.RunV("go", "build", "-o", "alb-responder", "./lambdas/alb-responder")
}

func ArchiveAlbResponder() error {
	mg.Deps(BuildAlbResponder)
	fmt.Println("Creating alb-responder archive...")

	return sh.RunV("zip", "alb-responder.zip", "alb-responder")
}

func LocalstackAvailable() error {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	start := time.Now()
	end := start.Add(time.Second * 60)
	for end.After(time.Now()) {
		fmt.Printf("Checking localstack availability at %s\n", time.Now().Sub(start))
		resp, err := client.Get("http://localstack:4566/health?reload")
		if err == nil {
			body, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			fmt.Printf("Got health response %s\n", string(body))
			if resp.StatusCode == 200 {
				fmt.Printf("Localstack available after %s\n", time.Now().Sub(start))
				return nil
			}
		} else {
			time.Sleep(time.Second * 5)
		}
	}
	return fmt.Errorf("timed out waiting for localstack")
}

func CreateAlbResponderBucket() error {
	mg.Deps(LocalstackAvailable)
	fmt.Println("Creating S3 bucket")
	return sh.RunV(
		"aws",
		"--endpoint-url=http://localstack:4566",
		"s3api",
		"create-bucket",
		"--bucket", "alb-responder",
		"--region", "eu-west-1",
	)
}

func Install() error {
	mg.Deps(LocalstackAvailable, ArchiveAlbResponder, CreateAlbResponderBucket)
	fmt.Println("Installing...")
	return sh.RunWithV(
		map[string]string{
			"CGO_ENABLED": "0",
		},
		"aws",
		"--endpoint-url=http://localstack:4566",
		"lambda",
		"create-function",
		"--role", "anything",
		"--function-name", "alb-responder",
		"--zip-file", "fileb://alb-responder.zip",
		"--environment", "Variables={AWS_ACCESS_KEY_ID=test,AWS_SECRET_ACCESS_KEY=test,AWS_DEFAULT_REGION=eu-west-1}",
		"--handler", "alb-responder",
		"--runtime", "go1.x",
	)
}

func DescribeLogs() error {
	mg.Deps(LocalstackAvailable)
	for i := 0; i < 20; i++ {
		sh.RunV(
			"aws",
			"--endpoint-url=http://localstack:4566",
			"logs",
			"describe-log-groups",
		)
		time.Sleep(time.Second * 10)
	}
	return nil
}

func PrintAlbResponderLogs() error {
	mg.Deps(LocalstackAvailable)
	for {
		sh.RunV(
			"aws",
			"--endpoint-url=http://localstack:4566",
			"logs", "get-log-events",
			"--log-group-name", "/aws/lambda/alb-responder",
			"--log-stream-name", "*",
		)
		time.Sleep(time.Second * 1)
	}
}
