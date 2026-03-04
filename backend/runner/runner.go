package runner

import (
	"bufio"
	"context"
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
	ctx context.Context
	mu  sync.Mutex
}

func NewCommandRunner(ctx context.Context) *CommandRunner {
	return &CommandRunner{ctx: ctx}
}

func (r *CommandRunner) RunCommand(cmdStr string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", cmdStr)
		setHideWindow(cmd)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go r.streamOutput(stdout)
	go r.streamOutput(stderr)

	return cmd.Wait()
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

		wailsRuntime.EventsEmit(r.ctx, "terminal:output", line)
	}
}

func (r *CommandRunner) RunRepairPhase(phase int) error {
	switch phase {
	case 1: // DISM
		r.Log("--- PHASE 1: DISM Health Check ---")
		steps := []string{
			"dism /online /cleanup-image /checkhealth",
			"dism /online /cleanup-image /scanhealth",
			"dism /online /cleanup-image /restorehealth",
		}
		for _, s := range steps {
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
		return r.RunCommand("sfc /scannow")
	case 4: // Component Cleanup
		r.Log("--- PHASE 4: Component Store Cleanup ---")
		return r.RunCommand("dism /online /cleanup-image /startcomponentcleanup /resetbase")
	case 5: // WMI
		r.Log("--- PHASE 5: WMI Repair ---")
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
			r.Log(fmt.Sprintf("> %s", s))
			if err := r.RunCommand(s); err != nil {
				// WMI repairs can be finicky; log error but try to continue for certain steps
				r.Log(fmt.Sprintf("! Warning: %s failed: %v", s, err))
			}
		}
	case 6: // AppX
		r.Log("--- PHASE 6: AppX Re-registration ---")
		return r.RunCommand("powershell -ExecutionPolicy Bypass -Command \"Get-AppXPackage -AllUsers | ForEach-Object { Add-AppxPackage -DisableDevelopmentMode -Register '$($_.InstallLocation)\\AppXManifest.xml' -ErrorAction SilentlyContinue }\"")
	default:
		return fmt.Errorf("invalid repair phase: %d", phase)
	}
	return nil
}

func (r *CommandRunner) Log(msg string) {
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
