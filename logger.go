package logger

import (
	"fmt"
	"sync"
	"time"
)

var globalMaxSiteRegionWidth int = 4
var globalMaxRowIndexWidth int = 1

var (
	logChan    = make(chan string, 100000)
	wg         sync.WaitGroup
	shutdownMu sync.Mutex
	shutdown   bool
)

const numWorkers = 4

func init() {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processLogMessages()
	}
}

func processLogMessages() {
	defer wg.Done()
	for msg := range logChan {
		fmt.Print(msg)
	}
}

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
	shutdownMu.Lock()
	defer shutdownMu.Unlock()

	if shutdown {
		return
	}
	timestamp := formatTime()
	message := fmt.Sprintf(format, args...)
	combined := l.site
	if l.region != "" {
		combined = l.site + " " + l.region
	}

	logEntry := fmt.Sprintf("%s[%-*s  %s  %*d] - %s%s%s\n",
		color,
		globalMaxSiteRegionWidth, combined,
		timestamp,
		globalMaxRowIndexWidth, l.rowIndex,
		message,
		ColorReset,
		"")

	logChan <- logEntry

	wg.Add(1)
	go func() {
		defer wg.Done()
		logChan <- logEntry
	}()
}

func Shutdown() {
	shutdownMu.Lock()
	shutdown = true
	shutdownMu.Unlock()

	close(logChan)
	wg.Wait()
}
