package main

import (
	"context"
	"fmt"
	"devale-v2/backend/runner"
	"devale-v2/backend/sysinfo"
	"devale-v2/backend/persistence"
)

// App struct
type App struct {
	ctx     context.Context
	runner  *runner.CommandRunner
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.runner = runner.NewCommandRunner(ctx)

	if !runner.IsAdmin() {
		a.runner.Log("! ERROR: DevAle is NOT running with Administrator privileges.")
		a.runner.Log("! Some system repair features may fail to execute correctly.")
		a.runner.Log("! Please restart the application as Administrator.")
	} else {
		a.runner.Log(">>> Administrator privileges verified.")
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) RunCommand(cmd string) string {
	// SANITIZATION: Check for malicious command injection patterns
	// in user-provided strings from the terminal
	restricted := []string{"&", "|", ";", ">", "<", "`", "$"}
	for _, char := range restricted {
		if strings.Contains(cmd, char) {
			return fmt.Sprintf("Error: Restricted character '%s' detected. Command blocked for security.", char)
		}
	}

	err := a.runner.RunCommand(cmd)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return "Success"
}

func (a *App) GetApplications() ([]runner.AppCategory, error) {
	return a.runner.GetApplications()
}

func (a *App) StopRepair() {
	a.runner.StopCommand()
}

func (a *App) GetSystemInfo() (*sysinfo.Info, error) {
	return sysinfo.GetSystemInfo()
}

func (a *App) SaveState(phase int) error {
	if phase < 0 || phase > 6 {
		return fmt.Errorf("invalid phase: %d", phase)
	}
	return persistence.SaveState(persistence.State{Phase: phase})
}

func (a *App) LoadState() (int, error) {
	state, err := persistence.LoadState()
	return state.Phase, err
}

func (a *App) RunRepairPhase(phase int) string {
	if phase < 1 || phase > 6 {
		return fmt.Sprintf("Error: invalid phase %d", phase)
	}
	err := a.runner.RunRepairPhase(phase)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return "Success"
}

func (a *App) ScheduleResume() string {
	err := a.runner.ScheduleResume()
	if err != nil {
		return err.Error()
	}
	return "Success"
}

func (a *App) ClearResume() string {
	err := a.runner.ClearResume()
	if err != nil {
		return err.Error()
	}
	return "Success"
}
