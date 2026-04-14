package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

func notifyDone(sessionType string) {
	title, message := notificationText(sessionType)

	switch runtime.GOOS {
	case "darwin":
		script := fmt.Sprintf(`display notification %q with title %q`, message, title)
		_ = exec.Command("osascript", "-e", script).Run()
	case "linux":
		// notify-send is available on most desktop Linux distros (libnotify).
		// Fall back to terminal bell if it's not installed.
		if err := exec.Command("notify-send", title, message).Run(); err != nil {
			fmt.Print("\a")
		}
	default:
		fmt.Print("\a")
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
