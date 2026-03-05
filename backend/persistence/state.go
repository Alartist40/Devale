package persistence

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type State struct {
	Phase int `json:"phase"`
}

func getStoragePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	path := filepath.Join(dir, "DevAle")
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}
	return filepath.Join(path, "state.json"), nil
}

func SaveState(state State) error {
	file, err := getStoragePath()
	if err != nil {
		return err
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}

func LoadState() (State, error) {
	file, err := getStoragePath()
	if err != nil {
		return State{}, err
	}
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
	file, err := getStoragePath()
	if err != nil {
		return err
	}
	return os.Remove(file)
}
