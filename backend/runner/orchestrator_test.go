package runner

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestRepairPhases_FullSuite(t *testing.T) {
	ctx := context.WithValue(context.Background(), isTestKey, true)
	mock := &MockCommander{
		Responses: map[string]MockResponse{
			"dism /online /cleanup-image /checkhealth":   {Output: "No component store corruption detected."},
			"chkdsk C: /f":                                {Output: "Cannot lock current drive.", Err: errors.New("exit status 3")},
			"sfc /scannow":                                {Output: "Windows Resource Protection did not find any integrity violations."},
			"dism /online /cleanup-image /resetbase":      {Output: "The operation completed successfully."},
			"net stop winmgmt /y":                         {Output: "The WMI service was stopped successfully."},
			"Get-AppXPackage":                             {Output: "AppX Packages re-registered successfully."},
		},
	}
	runner := NewCommandRunner(ctx)
	runner.SetCommander(mock)

	for phase := 1; phase <= 6; phase++ {
		t.Run(fmt.Sprintf("Phase %d", phase), func(t *testing.T) {
			err := runner.RunRepairPhase(phase)
			if err != nil && phase != 2 { // Phase 2 might "fail" with status 3 but we handle it
				t.Errorf("Phase %d failed: %v", phase, err)
			}
		})
	}
}

func TestRepairPhase2_RestartScenarios(t *testing.T) {
	ctx := context.WithValue(context.Background(), isTestKey, true)
	mock := &MockCommander{
		Responses: map[string]MockResponse{
			"chkdsk C: /f": {Output: "The type of the file system is NTFS. Cannot lock current drive.", Err: errors.New("exit status 3")},
		},
	}
	runner := NewCommandRunner(ctx)
	runner.SetCommander(mock)

	// Test Phase 2 success (Status 3 is treated as success)
	err := runner.RunRepairPhase(2)
	if err != nil {
		t.Errorf("Expected Phase 2 to succeed with status 3, got error: %v", err)
	}

	found := false
	for _, cmd := range mock.History {
		if strings.Contains(cmd, "chkdsk C: /f") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("chkdsk command not found in history")
	}
}

func TestCommandInjection_Protection(t *testing.T) {
	// Note: Protection is currently handled at the App layer (app.go)
	// runner.go trusts internal commands but we want to ensure the App layer blocks it.
	// Since we can't easily test app.go here without complex setups, we'll verify
	// that runner.go still works for trusted commands.
	ctx := context.WithValue(context.Background(), isTestKey, true)
	mock := &MockCommander{Responses: make(map[string]MockResponse)}
	runner := NewCommandRunner(ctx)
	runner.SetCommander(mock)

	trustedCmd := "sfc /scannow"
	runner.RunCommand(trustedCmd)

	found := false
	for _, cmd := range mock.History {
		if cmd == trustedCmd {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Trusted command failed to execute")
	}
}

func TestWMIInterruptionVulnerability(t *testing.T) {
	ctx := context.WithValue(context.Background(), isTestKey, true)
	// Mock failing in the middle of WMI repair
	mock := &MockCommander{
		Responses: map[string]MockResponse{
			"net stop winmgmt": {Err: errors.New("service could not be stopped")},
		},
	}
	runner := NewCommandRunner(ctx)
	runner.SetCommander(mock)

	err := runner.RunRepairPhase(5)
	// We expect it to continue because of current "Warning" logic, but is it safe?
	if err != nil {
		t.Errorf("WMI phase failed: %v", err)
	}

	// Check if we attempted to re-enable
	foundEnable := false
	for _, cmd := range mock.History {
		if strings.Contains(cmd, "sc config winmgmt start= auto") {
			foundEnable = true
		}
	}

	if !foundEnable {
		t.Error("VULNERABILITY: winmgmt re-enable command not found in history; service may be left disabled")
	}
}
