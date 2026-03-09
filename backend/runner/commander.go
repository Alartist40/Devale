package runner

import (
	"context"
	"fmt"
	"io"
	"strings"
)

// Commander interface allows us to mock execution
type Commander interface {
	Run(ctx context.Context, cmd string, stdout, stderr io.Writer) error
}

// RealCommander runs actual system commands (Implemented in runner.go)
type RealCommander struct{}

// MockCommander simulates Windows responses
type MockCommander struct {
	Responses map[string]MockResponse
	History   []string
}

type MockResponse struct {
	Output string
	Err    error
}

func (c *MockCommander) Run(ctx context.Context, cmdStr string, stdout, stderr io.Writer) error {
	c.History = append(c.History, cmdStr)

	for pattern, resp := range c.Responses {
		if strings.Contains(cmdStr, pattern) {
			if resp.Output != "" {
				fmt.Fprintln(stdout, resp.Output)
			}
			return resp.Err
		}
	}

	// Default success for unknown commands
	return nil
}
