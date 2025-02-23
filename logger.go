package logger

import (
	"fmt"
	"sync"
	"time"
)

var globalMaxSiteRegionWidth int = 4
var globalMaxRowIndexWidth int = 1

var (
	logChan     chan string
	wgWorkers   sync.WaitGroup
	wgProducers sync.WaitGroup
	shutdownMu  sync.Mutex
	shutdown    bool
	numWorkers  = 4
	initialized bool
	initMu      sync.Mutex
)

type Logger struct {
	site     string
	region   string
	rowIndex int
}

const (
	ColorGreen  = "\033[32m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorPurple = "\033[35m"
	ColorReset  = "\033[0m"
)

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

func InitializeLogger() {
	initMu.Lock()
	defer initMu.Unlock()

	if initialized {
		return
	}

	logChan = make(chan string, 100000)

	// Start worker pool
	for i := 0; i < numWorkers; i++ {
		wgWorkers.Add(1)
		go func(workerID int) {
			defer wgWorkers.Done()
			for msg := range logChan {
				fmt.Print(msg)
			}
		}(i)
	}

	initialized = true
	shutdown = false
}

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
	if shutdown {
		shutdownMu.Unlock()
		return
	}
	shutdownMu.Unlock()

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

	wgProducers.Add(1)
	go func(entry string) {
		defer wgProducers.Done()
		logChan <- entry
	}(logEntry)
}

func startLogProcessor(workerCount int) {
	for i := 0; i < workerCount; i++ {
		wgWorkers.Add(1)
		go func(workerID int) {
			defer wgWorkers.Done()
			for msg := range logChan {
				fmt.Print(msg)
			}
		}(i)
	}
}

func Shutdown() {
	shutdownMu.Lock()
	if shutdown {
		shutdownMu.Unlock()
		return
	}
	shutdown = true
	shutdownMu.Unlock()

	wgProducers.Wait()

	close(logChan)

	wgWorkers.Wait()
	initMu.Lock()
	defer initMu.Unlock()
	initialized = false
}

func init() {
	startLogProcessor(numWorkers)
}
