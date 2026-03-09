package persistence

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type State struct {
	Phase      int      `json:"phase"`
	LastStep   string   `json:"last_step,omitempty"`
	Logs       []string `json:"logs,omitempty"`
	IsComplete bool     `json:"is_complete"`
}

func getStoragePath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	path := filepath.Join(dir, "DevAle")
	os.MkdirAll(path, 0755)
	return filepath.Join(path, "state.json")
}

func SaveState(state State) error {
	file := getStoragePath()
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}

func LoadState() (State, error) {
	file := getStoragePath()
	var state State
	data, err := os.ReadFile(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return state, nil
		}
		return state, err
	}

	err = json.Unmarshal(data, &state)
	return state, err
}

func ClearState() error {
	file := getStoragePath()
	return os.Remove(file)
}
