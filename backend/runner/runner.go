package runner

import (
	"bufio"
	"context"
	"devale-v2/backend/persistence"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/text/encoding/charmap"
)

type Step struct {
	Name    string
	Command string
}

type CommandRunner struct {
	ctx        context.Context
	mu         sync.Mutex
	cancelMu   sync.Mutex
	cancelFunc context.CancelFunc
	commander  Commander
}

func NewCommandRunner(ctx context.Context) *CommandRunner {
	return &CommandRunner{
		ctx:       ctx,
		commander: &RealCommander{},
	}
}

// SetCommander allows swapping the execution engine (e.g. for testing)
func (r *CommandRunner) SetCommander(c Commander) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commander = c
}

func (r *CommandRunner) RunCommand(cmdStr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create a sub-context for this command
	cmdCtx, cancel := context.WithCancel(r.ctx)
	r.cancelMu.Lock()
	r.cancelFunc = cancel
	r.cancelMu.Unlock()

	defer func() {
		r.cancelMu.Lock()
		r.cancelFunc = nil
		r.cancelMu.Unlock()
	}()

	pr, pw := io.Pipe()
	go r.streamOutput(pr)

	err := r.commander.Run(cmdCtx, cmdStr, pw, pw)
	pw.Close()
	return err
}

func (c *RealCommander) Run(ctx context.Context, cmdStr string, stdout, stderr io.Writer) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/c", cmdStr)
		setHideWindow(cmd)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", cmdStr)
	}

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func (r *CommandRunner) streamOutput(reader io.ReadCloser) {
	defer reader.Close()

	var rdr io.Reader = reader
	if runtime.GOOS == "windows" {
		rdr = charmap.CodePage850.NewDecoder().Reader(reader)
	}

	scanner := bufio.NewScanner(rdr)
	for scanner.Scan() {
		line := scanner.Text()

		// Filter out Windows shell noise
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Skip copyright and version headers
		if strings.Contains(line, "Microsoft Windows [Version") ||
			strings.Contains(line, "(c) Microsoft Corporation") ||
			strings.Contains(line, "Microsoft (R) Windows (R)") {
			continue
		}

		// FILTERING FOR CLEANER UI (No glitching/scrolling noise)
		// Skip verbose per-file logs for common long-running operations
		if strings.Contains(line, "regsvr32 /s") ||
			strings.Contains(line, "mofcomp ") ||
			strings.Contains(line, "Parsing MOF file:") ||
			strings.Contains(line, "MOF file has been successfully parsed") ||
			strings.Contains(line, "Storing data in the repository") ||
			strings.Contains(line, "Verification") && strings.Contains(line, "% complete") ||
			strings.Contains(line, "[") && strings.Contains(line, "%") && strings.Contains(line, "=") {

			// Periodically send a "Working..." message instead of every line if we want to show life
			// For now, let's just emit the phase header from the orchestrator and keep terminal clean
			continue
		}

		r.Log(line)
	}
}

func (r *CommandRunner) StopCommand() {
	r.cancelMu.Lock()
	defer r.cancelMu.Unlock()
	if r.cancelFunc != nil {
		r.cancelFunc()
		r.Log("!!! PROCESS MANUALLY INTERRUPTED BY USER !!!")
	}
}

func (r *CommandRunner) RunRepairPhase(phase int) error {
	select {
	case <-r.ctx.Done():
		return fmt.Errorf("operation cancelled")
	default:
	}

	updateLastStep := func(step string) {
		state, err := persistence.LoadState()
		if err == nil {
			state.Phase = phase
			state.LastStep = step
			persistence.SaveState(state)
		}
	}

	switch phase {
	case 1: // DISM
		r.Log("--- PHASE 1: DISM Health Check ---")
		steps := []string{
			"dism /online /cleanup-image /checkhealth",
			"dism /online /cleanup-image /scanhealth",
			"dism /online /cleanup-image /restorehealth",
		}
		for _, s := range steps {
			updateLastStep(s)
			if err := r.RunCommand(s); err != nil {
				return fmt.Errorf("phase 1 step failed (%s): %w", s, err)
			}
		}
	case 2: // CHKDSK
		r.Log("--- PHASE 2: Scheduling CHKDSK ---")
		err := r.RunCommand("echo Y | chkdsk C: /f")
		if err != nil {
			// Exit status 3 is common when chkdsk schedules successfully but cannot lock the drive
			if strings.Contains(err.Error(), "exit status 3") {
				r.Log(">>> CHKDSK scheduled successfully (status 3).")
				return nil
			}
			return err
		}
		return nil
	case 3: // SFC
		r.Log("--- PHASE 3: SFC Scan ---")
		r.Log(">>> Beginning system scan. This will take 10-15 minutes...")
		return r.RunCommand("sfc /scannow")
	case 4: // Component Cleanup
		r.Log("--- PHASE 4: Component Store Cleanup ---")
		r.Log(">>> Cleaning up component store. This will take several minutes...")
		return r.RunCommand("dism /online /cleanup-image /startcomponentcleanup /resetbase")
	case 5: // WMI
		r.Log("--- PHASE 5: WMI Repair ---")
		r.Log(">>> Stopping WMI services and re-registering core libraries...")

		// Ensure we re-enable WMI no matter what happens
		defer func() {
			r.Log(">>> Ensuring WMI service is re-enabled...")
			r.RunCommand("sc config winmgmt start= auto")
			r.RunCommand("net start winmgmt")
		}()

		steps := []string{
			"sc config winmgmt start= disabled",
			"net stop winmgmt /y",
			"cd /d %windir%\\system32\\wbem && for /f %s in ('dir /b *.dll') do regsvr32 /s %s",
			"wmiprvse /regserver",
			"winmgmt /regserver",
			"sc config winmgmt start= auto",
			"net start winmgmt",
			"cd /d %windir%\\system32\\wbem && for /f %s in ('dir /b *.mof *.mfl ^| findstr /v /i \"uninstall\"') do mofcomp %s",
		}
		for _, s := range steps {
			updateLastStep(s)
			r.Log(fmt.Sprintf("> %s", s))
			if err := r.RunCommand(s); err != nil {
				// WMI repairs can be finicky; log error but try to continue for certain steps
				r.Log(fmt.Sprintf("! Warning: %s failed: %v", s, err))
			}

			// Check if we were cancelled
			select {
			case <-r.ctx.Done():
				return fmt.Errorf("WMI phase interrupted")
			default:
			}
		}
	case 6: // AppX
		r.Log("--- PHASE 6: AppX Re-registration ---")
		return r.RunCommand("powershell -ExecutionPolicy Bypass -Command \"Get-AppXPackage -AllUsers | ForEach-Object { Add-AppxPackage -DisableDevelopmentMode -Register \\\"$($_.InstallLocation)\\AppXManifest.xml\\\" -ErrorAction SilentlyContinue }\"")
	default:
		return fmt.Errorf("invalid repair phase: %d", phase)
	}
	return nil
}

type contextKey string

const isTestKey contextKey = "isTest"

func (r *CommandRunner) Log(msg string) {
	if r.ctx == nil {
		fmt.Printf("[LOG] %s\n", msg)
		return
	}

	// Check for our custom test key
	if val := r.ctx.Value(isTestKey); val != nil && val.(bool) {
		fmt.Printf("[LOG] %s\n", msg)
		return
	}

	// Safely attempt to emit, but don't crash
	defer func() { recover() }()
	wailsRuntime.EventsEmit(r.ctx, "terminal:output", msg)
}

func (r *CommandRunner) ScheduleResume() error {
	if runtime.GOOS != "windows" {
		return nil
	}
	exe, err := os.Executable()
	if err != nil {
		exe = "devale-v2.exe"
	}
	cmd := fmt.Sprintf("schtasks /create /tn \"DevAleResume\" /tr \"\\\"%s\\\"\" /sc onlogon /rl highest /f", exe)
	return r.RunCommand(cmd)
}

func (r *CommandRunner) ClearResume() error {
	if runtime.GOOS != "windows" {
		return nil
	}
	return r.RunCommand("schtasks /delete /tn \"DevAleResume\" /f")
}

func (r *CommandRunner) GetApplications() ([]AppCategory, error) {
	// Attempt to load from relative path
	data, err := os.ReadFile("backend/runner/applications.json")
	if err != nil {
		// Fallback for build environment
		data, err = os.ReadFile("applications.json")
		if err != nil {
			return nil, fmt.Errorf("could not load applications.json: %w", err)
		}
	}

	var result struct {
		Categories []AppCategory `json:"categories"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result.Categories, nil
}
