package logger

import (
	"fmt"
	"time"
)

type Logger struct {
	site     string
	region   string
	rowIndex int
}

func Setup(site, region string, rowIndex int) *Logger {
	return &Logger{
		site:     site,
		region:   region,
		rowIndex: rowIndex,
	}
}

func formatTime() string {
	return time.Now().Format("15:04:05.00")
}

const (
	ColorGreen  = "\033[32m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorPurple = "\033[35m"
	ColorReset  = "\033[0m"
)

// green
func (l *Logger) Success(msg string) {
	timestamp := formatTime()
	fmt.Printf("[%s %s  %s  %d] - %s%s%s\n", l.site, l.region, timestamp, l.rowIndex, ColorGreen, msg, ColorReset)
}

// red
func (l *Logger) Error(msg string) {
	timestamp := formatTime()
	fmt.Printf("[%s %s  %s  %d] - %s%s%s\n", l.site, l.region, timestamp, l.rowIndex, ColorRed, msg, ColorReset)
}

// yellow
func (l *Logger) Info(msg string) {
	timestamp := formatTime()
	fmt.Printf("[%s %s  %s  %d] - %s%s%s\n", l.site, l.region, timestamp, l.rowIndex, ColorYellow, msg, ColorReset)
}

// purple
func (l *Logger) Normal(msg string) {
	timestamp := formatTime()
	fmt.Printf("[%s %s  %s  %d] - %s%s%s\n", l.site, l.region, timestamp, l.rowIndex, ColorPurple, msg, ColorReset)
}
