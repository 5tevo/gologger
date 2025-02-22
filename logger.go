package logger

import (
	"fmt"
	"time"
)

var globalMaxSiteRegionWidth int = 4
var globalMaxRowIndexWidth int = 1

type Logger struct {
	site     string
	region   string
	rowIndex int
}

func SetGlobalWidths(siteRegionWidth, rowIndexWidth int) {
	globalMaxSiteRegionWidth = siteRegionWidth
	globalMaxRowIndexWidth = rowIndexWidth
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
	l.logMessage(ColorGreen, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.logMessage(ColorRed, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.logMessage(ColorYellow, format, args...)
}

func (l *Logger) Normal(format string, args ...interface{}) {
	l.logMessage(ColorPurple, format, args...)
}

func (l *Logger) logMessage(color, format string, args ...interface{}) {
	timestamp := formatTime()
	message := fmt.Sprintf(format, args...)
	combined := l.site
	if l.region != "" {
		combined = l.site + " " + l.region
	}

	fmt.Printf("%s[%-*s  %s  %*d] - %s%s%s\n",
		color,
		globalMaxSiteRegionWidth, combined,
		timestamp,
		globalMaxRowIndexWidth, l.rowIndex,
		message,
		ColorReset,
		"",
	)
}
