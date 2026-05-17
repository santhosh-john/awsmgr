package ui

import (
	"fmt"

	"github.com/santhoshjohn/awsmgr/internal/doctor"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

func PrintDoctorReport(report doctor.Report) {
	fmt.Println("Environment diagnostics:")
	fmt.Println()
	for _, check := range report.Checks {
		PrintStatus(string(check.ResultStatus()), check.Name, check.Message)
	}
}

func PrintStatus(status, area, message string) {
	var color string
	switch status {
	case "OK":
		color = ColorGreen
	case "FAIL":
		color = ColorRed
	case "WARN":
		color = ColorYellow
	default:
		color = ColorReset
	}

	fmt.Printf("[%s%-4s%s]  %s%-15s%s %s\n",
		color, status, ColorReset,
		ColorBold, area, ColorReset,
		message,
	)
}
