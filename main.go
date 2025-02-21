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

func (l *Logger) Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	timestamp := formatTime()
	fmt.Printf("[%s %s  %s  %d] - %s%s%s\n", l.site, l.region, timestamp, l.rowIndex, ColorGreen, msg, ColorReset)
}

func (l *Logger) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	timestamp := formatTime()
	fmt.Printf("[%s %s  %s  %d] - %s%s%s\n", l.site, l.region, timestamp, l.rowIndex, ColorRed, msg, ColorReset)
}

func (l *Logger) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	timestamp := formatTime()
	fmt.Printf("[%s %s  %s  %d] - %s%s%s\n", l.site, l.region, timestamp, l.rowIndex, ColorYellow, msg, ColorReset)
}

func (l *Logger) Normal(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	timestamp := formatTime()
	fmt.Printf("[%s %s  %s  %d] - %s%s%s\n", l.site, l.region, timestamp, l.rowIndex, ColorPurple, msg, ColorReset)
}
