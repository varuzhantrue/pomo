package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

func notifyDone(sessionType string) {
	// Terminal bell — works everywhere.
	fmt.Print("\a")

	// Native macOS notification via osascript.
	if runtime.GOOS == "darwin" {
		title, message := notificationText(sessionType)
		script := fmt.Sprintf(`display notification %q with title %q`, message, title)
		// Ignore errors — notification is best-effort.
		_ = exec.Command("osascript", "-e", script).Run()
	}
}

func notificationText(sessionType string) (title, message string) {
	switch sessionType {
	case "focus":
		return "Pomodoro Complete! 🍅", "Time for a break. Well done!"
	case "break":
		return "Break Over! ☕", "Time to focus again."
	case "long":
		return "Long Break Over! 🌿", "Ready to get back to work?"
	default:
		return "Pomo", "Session complete."
	}
}
