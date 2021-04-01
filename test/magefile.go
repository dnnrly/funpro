// +build mage

package main

import (
	"fmt"
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
	end := start.Add(time.Second * 30)
	for end.After(time.Now()) {
		fmt.Printf("Checking localstack availability at %s\n", time.Now().Sub(start))
		resp, err := client.Get("http://localstack:4566/health?reload")
		if err == nil {
			defer resp.Body.Close()
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
	mg.Deps(ArchiveAlbResponder, CreateAlbResponderBucket)
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
		"--handler", "main",
		"--runtime", "go1.x",
	)
}
