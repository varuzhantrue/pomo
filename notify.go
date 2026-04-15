package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

//go:embed assets/pomo-notification.ogg
var notificationSound []byte

func notifyDone(sessionType string) {
	title, message := notificationText(sessionType)

	switch runtime.GOOS {
	case "darwin":
		script := fmt.Sprintf(`display notification %q with title %q`, message, title)
		_ = exec.Command("osascript", "-e", script).Run()
		playSound("afplay")
	case "linux":
		// notify-send is available on most desktop Linux distros (libnotify).
		// Fall back to terminal bell if it's not installed.
		if err := exec.Command("notify-send", title, message).Run(); err != nil {
			fmt.Print("\a")
		}
		playSound("paplay")
	default:
		fmt.Print("\a")
	}
}

func playSound(player string) {
	if loadConfig().Muted {
		return
	}
	if _, err := exec.LookPath(player); err != nil {
		return // player not installed
	}
	f, err := os.CreateTemp("", "pomo-*.ogg")
	if err != nil {
		return
	}
	defer os.Remove(f.Name())

	if _, err := f.Write(notificationSound); err != nil {
		f.Close()
		return
	}
	f.Close()

	_ = exec.Command(player, f.Name()).Run()
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
