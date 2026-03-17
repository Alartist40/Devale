package runner

import (
	"context"
	"devale-v2/backend/persistence"
	"testing"
	"time"
)

func TestPhase5Cancellation(t *testing.T) {
	ctx := context.WithValue(context.Background(), isTestKey, true)
	runner := NewCommandRunner(ctx)

	// Reset persistence to ensure we are starting fresh
	persistence.SaveState(persistence.State{})

	mock := &MockCommander{
		Responses: map[string]MockResponse{
			"sc config winmgmt start= disabled": {Output: "Success"},
			"net stop winmgmt /y":                {Output: "Stopped"},
		},
		// We want to simulate a cancellation during the third step
		Delay: map[string]time.Duration{
			"cd /d %windir%\\system32\\wbem && for /f %s in ('dir /b *.dll') do regsvr32 /s %s": 100 * time.Millisecond,
		},
	}
	runner.SetCommander(mock)

	// Start Phase 5 in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- runner.RunRepairPhase(5)
	}()

	// Wait a bit and then cancel
	time.Sleep(50 * time.Millisecond)
	runner.StopCommand()

	// Wait for Phase 5 to finish
	select {
	case err := <-errChan:
		if err == nil {
			t.Error("Expected error from cancelled Phase 5, got nil")
		}
		// Check if it continued to subsequent steps
		for _, cmd := range mock.History {
			if cmd == "wmiprvse /regserver" {
				t.Error("VULNERABILITY: Phase 5 continued to 'wmiprvse /regserver' after cancellation")
			}
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for Phase 5 to respond to cancellation")
	}
}
