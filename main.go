package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	focusDuration = 25 * 60 // seconds
	breakDuration = 5 * 60
	longDuration  = 15 * 60
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		n := len(todaySessions()) + 1
		RunTimer("focus", focusDuration, n)

	case "break":
		RunTimer("break", breakDuration, 0)

	case "long":
		RunTimer("long", longDuration, 0)

	case "status":
		showStatus()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %q\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func showStatus() {
	sessions := todaySessions()
	if len(sessions) == 0 {
		fmt.Println("No focus sessions yet today. Run `pomo start` to begin!")
		return
	}

	total := 0
	for _, s := range sessions {
		total += s.Duration
	}

	tomatoes := strings.Repeat("🍅", len(sessions))
	fmt.Printf("Today: %s  (%d sessions = %s focused)\n",
		tomatoes, len(sessions), focusTimeString(total))

	last := sessions[len(sessions)-1]
	fmt.Printf("Last session started: %s\n", last.Timestamp.Format(time.Kitchen))
}

func printUsage() {
	fmt.Print(`pomo — a Pomodoro timer for your terminal

Usage:
  pomo start    25-minute focus session
  pomo break     5-minute short break
  pomo long     15-minute long break
  pomo status   show today's sessions and focus time
`)
}
