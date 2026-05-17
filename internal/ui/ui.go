package ui

import (
	"fmt"
	"os"

	"github.com/santhosh-john/awsmgr/internal/doctor"
)

const (
	ColorReset      = "\033[0m"
	ColorRed        = "\033[31m"
	ColorGreen      = "\033[32m"
	ColorYellow     = "\033[33m"
	ColorCyan       = "\033[36m"
	ColorBold       = "\033[1m"
	ColorBoldRed    = "\033[1;31m"
	ColorBoldGreen  = "\033[1;32m"
	ColorBoldYellow = "\033[1;33m"
	ColorBoldCyan   = "\033[1;36m"
)

// colorsEnabled checks if colors should be output to terminal.
// Returns false if NO_COLOR env var is set or output is not a TTY.
func ColorsEnabled() bool {
	// Respect NO_COLOR environment variable (https://no-color.org/)
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check if stdout is a terminal
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// IsTerminal reports whether stdout is connected to a terminal.
func IsTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// ApplyColor wraps text with ANSI color codes if colors are enabled.
func ApplyColor(color, text string) string {
	if !ColorsEnabled() {
		return text
	}
	return color + text + ColorReset
}

func PrintDoctorReport(report doctor.Report) {
	fmt.Println("Environment diagnostics:")
	fmt.Println()
	for _, check := range report.Checks {
		PrintStatus(string(check.ResultStatus()), check.Name, check.Message)
	}
}

func PrintStatus(status, area, message string) {
	statusColor := statusColor(status)

	// Format status badge and area with conditional coloring
	var statusStr, areaStr string
	if ColorsEnabled() {
		statusStr = fmt.Sprintf("[%s%-4s%s]", statusColor, status, ColorReset)
		areaStr = fmt.Sprintf("%s%-15s%s", ColorBold, area, ColorReset)
	} else {
		statusStr = fmt.Sprintf("[%-4s]", status)
		areaStr = fmt.Sprintf("%-15s", area)
	}

	fmt.Printf("%s  %s %s\n", statusStr, areaStr, message)
}

// statusColor returns the appropriate ANSI color code for a status.
func statusColor(status string) string {
	switch status {
	case "OK":
		return ColorGreen
	case "FAIL":
		return ColorRed
	case "WARN":
		return ColorYellow
	default:
		return ColorReset
	}
}
