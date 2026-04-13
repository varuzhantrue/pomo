package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Session represents a single pomodoro session entry.
type Session struct {
	Type      string    `json:"type"`      // "focus", "break", "long"
	Timestamp time.Time `json:"timestamp"` // when the session started
	Duration  int       `json:"duration"`  // elapsed seconds
	Partial   bool      `json:"partial,omitempty"`
}

func dataFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".pomo")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "sessions.json"), nil
}

func loadSessions() ([]Session, error) {
	path, err := dataFile()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var sessions []Session
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func saveSession(s Session) error {
	sessions, err := loadSessions()
	if err != nil {
		// Don't lose data on a read error — start fresh.
		sessions = nil
	}
	sessions = append(sessions, s)

	path, err := dataFile()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// todaySessions returns completed (non-partial) focus sessions from today.
func todaySessions() []Session {
	sessions, err := loadSessions()
	if err != nil {
		return nil
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var result []Session
	for _, s := range sessions {
		if s.Type == "focus" && !s.Partial && s.Timestamp.After(today) {
			result = append(result, s)
		}
	}
	return result
}
