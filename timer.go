package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const barWidth = 20

// ANSI codes.
const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	cursorHide  = "\033[?25l"
	cursorShow  = "\033[?25h"
	clearLine   = "\r\033[K"
)

func pickColor(remainingSecs int) string {
	switch {
	case remainingSecs < 60:
		return colorRed
	case remainingSecs < 300:
		return colorYellow
	default:
		return colorGreen
	}
}

func renderBar(elapsed, total int) string {
	filled := int(float64(elapsed) / float64(total) * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled) + "]"
}

func fmtDuration(secs int) string {
	if secs < 0 {
		secs = 0
	}
	return fmt.Sprintf("%02d:%02d", secs/60, secs%60)
}

func sessionEmoji(t string) string {
	switch t {
	case "focus":
		return "🍅"
	case "break":
		return "☕"
	case "long":
		return "🌿"
	}
	return "⏱"
}

func sessionLabel(t string) string {
	switch t {
	case "focus":
		return "Focus Session"
	case "break":
		return "Short Break"
	case "long":
		return "Long Break"
	}
	return "Session"
}

// RunTimer runs a countdown for totalSecs seconds, displaying a live progress bar.
// sessionNum is used only for focus sessions to display the session number.
func RunTimer(sessionType string, totalSecs int, sessionNum int) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	start := time.Now()

	// Build the session header string once.
	var header string
	if sessionType == "focus" {
		header = fmt.Sprintf("%s %s #%d", sessionEmoji(sessionType), sessionLabel(sessionType), sessionNum)
	} else {
		header = fmt.Sprintf("%s %s", sessionEmoji(sessionType), sessionLabel(sessionType))
	}

	render := func(elapsed int) {
		remaining := totalSecs - elapsed
		c := pickColor(remaining)
		bar := renderBar(elapsed, totalSecs)
		fmt.Printf("%s%s%s  %s  %s remaining%s",
			clearLine, c, header, bar, fmtDuration(remaining), colorReset)
	}

	// Hide cursor for cleaner rendering; restore on any exit path.
	fmt.Print(cursorHide)
	defer fmt.Print(cursorShow)

	render(0)

	interrupted := false
	elapsed := 0

loop:
	for {
		select {
		case <-sigCh:
			interrupted = true
			elapsed = int(time.Since(start).Seconds())
			fmt.Println() // move past the progress line
			break loop

		case <-ticker.C:
			elapsed = int(time.Since(start).Seconds())
			render(elapsed)
			if elapsed >= totalSecs {
				fmt.Println()
				break loop
			}
		}
	}

	if interrupted {
		// Cursor is still hidden here; show it so the prompt is usable.
		fmt.Print(cursorShow)

		if elapsed < 1 {
			fmt.Println("Session cancelled.")
			return
		}

		fmt.Printf("Interrupted after %s. Save partial session? (y/n): ", fmtDuration(elapsed))
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))

		if answer == "y" || answer == "yes" {
			if err := saveSession(Session{
				Type:      sessionType,
				Timestamp: start,
				Duration:  elapsed,
				Partial:   true,
			}); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to save session: %v\n", err)
			} else {
				fmt.Println("Partial session saved.")
			}
		} else {
			fmt.Println("Session discarded.")
		}
		return
	}

	// --- Session completed successfully ---
	notifyDone(sessionType)

	if err := saveSession(Session{
		Type:      sessionType,
		Timestamp: start,
		Duration:  totalSecs,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save session: %v\n", err)
	}

	printCompletionMessage(sessionType)
}

func printCompletionMessage(sessionType string) {
	switch sessionType {
	case "focus":
		sessions := todaySessions()
		count := len(sessions)
		total := 0
		for _, s := range sessions {
			total += s.Duration
		}
		tomatoes := strings.Repeat("🍅", count)
		timeStr := focusTimeString(total)
		fmt.Printf("%s✅ Done! Take a break.%s\n", colorBold, colorReset)
		fmt.Printf("Today: %s  (%d sessions = %s focused)\n", tomatoes, count, timeStr)

	case "break":
		fmt.Printf("%s✅ Break over! Time to focus.%s\n", colorBold, colorReset)

	case "long":
		fmt.Printf("%s✅ Long break over! Ready to go?%s\n", colorBold, colorReset)
	}
}

func focusTimeString(totalSecs int) string {
	h := totalSecs / 3600
	m := (totalSecs % 3600) / 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}
