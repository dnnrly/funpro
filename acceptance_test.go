package main_test

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

type testContext struct {
	err       error
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
		Timeout: time.Second * 5,
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

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	var tc testContext
	ctx.BeforeScenario(func(*godog.Scenario) {
	})
	ctx.Step(`^the app runs with parameters "([^"]*)"$`, tc.theAppRunsWithParameters)
	ctx.Step(`^the app exits without error$`, tc.theAppExitsWithoutError)
	ctx.Step(`^the output contains "([^"]*)"$`, tc.theOutputContains)
	ctx.Step(`^I make a "([^"]*)" to "([^"]*)"$`, tc.iMakeARequestTo)
	ctx.Step(`^the response code is (\d+)$`, tc.theResponseCodeIs)
}
