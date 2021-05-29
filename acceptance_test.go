package main_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

type testContext struct {
	err      error
	s3Client *s3.S3

	cmdResult struct {
		Output string
		Err    error
	}
	response struct {
		code int
	}
}

// Errorf is used by the called assertion to report an error and is required to
// make testify assertions work
func (c *testContext) Errorf(format string, args ...interface{}) {
	c.err = fmt.Errorf(format, args...)
}

func (c *testContext) theAppRunsWithParameters(args string) error {
	cmd := exec.Command("./funpro", strings.Split(args, " ")...)
	output, err := cmd.CombinedOutput()
	c.cmdResult.Output = string(output)
	c.cmdResult.Err = err

	return nil
}

func (c *testContext) theAppExitsWithoutError() error {
	assert.NoError(c, c.cmdResult.Err)
	return c.err
}

func (c *testContext) iMakeARequestTo(method, url string) error {
	client := http.Client{
		Timeout: time.Second * 360,
	}
	body := strings.NewReader("")
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	c.response.code = resp.StatusCode
	return nil
}

func (c *testContext) theResponseCodeIs(expected int) error {
	assert.Equal(c, expected, c.response.code)
	return c.err
}

func (c *testContext) theOutputContains(expected string) error {
	assert.Contains(c, c.cmdResult.Output, expected)
	return c.err
}

func (c *testContext) theRecordingAtMatches(bucket, path, expectedPath string) error {
	req, resp := c.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})

	err := req.Send()
	if err != nil {
		return fmt.Errorf("Cannot fetch %s:%s %w", bucket, path, err)
	}

	expectedFile, err := os.Open("test/data/" + expectedPath)
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	expected, err := io.ReadAll(expectedFile)
	if err != nil {
		return err
	}

	assert.JSONEq(c, string(expected), string(buf))

	return nil
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	awsConfig := aws.NewConfig().
		WithRegion("eu-west-1").
		WithEndpoint("http://localstack:4566")
	sess := session.Must(session.NewSession(awsConfig))

	tc := testContext{
		s3Client: s3.New(sess, awsConfig.WithS3ForcePathStyle(true)),
	}
	ctx.BeforeScenario(func(*godog.Scenario) {
	})
	ctx.AfterScenario(func(sc *godog.Scenario, err error) {
		if err != nil {
			fmt.Printf("Exit code %d\nOutput:\n%s\n", tc.response.code, tc.cmdResult.Output)
		}
	})
	ctx.Step(`^the app runs with parameters "([^"]*)"$`, tc.theAppRunsWithParameters)
	ctx.Step(`^the app exits without error$`, tc.theAppExitsWithoutError)
	ctx.Step(`^the output contains "([^"]*)"$`, tc.theOutputContains)
	ctx.Step(`^I make a "([^"]*)" to "([^"]*)"$`, tc.iMakeARequestTo)
	ctx.Step(`^the response code is (\d+)$`, tc.theResponseCodeIs)
	ctx.Step(`^the "([^"]*)" recording at "([^"]*)" matches "([^"]*)"$`, tc.theRecordingAtMatches)
}
