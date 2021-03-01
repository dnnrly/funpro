package main_test

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

type testContext struct {
	err    error
	result struct {
		Output string
		Err    error
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
	c.result.Output = string(output)
	c.result.Err = err

	return nil
}

func (c *testContext) theAppExitsWithoutError() error {
	assert.NoError(c, c.result.Err)
	return c.err
}

func (c *testContext) theOutputContains(expected string) error {
	assert.Contains(c, c.result.Output, expected)
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
}
