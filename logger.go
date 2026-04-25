package logger

import (
	"fmt"
	"sync"
	"time"
)

var (
	globalMaxSiteRegionWidth int = 4
	globalMaxRowIndexWidth   int = 1
	logChan                  chan string
	wgWorkers                sync.WaitGroup
	wgProducers              sync.WaitGroup
	shutdownMu               sync.Mutex
	shutdown                 bool
	numWorkers               = 4
	initialised              bool
	initMu                   sync.Mutex
	channelSize              = 100000
)

type Logger struct {
	combined string
	rowIndex int
}

const (
	ColorGreen     = "\033[32m"
	ColorRed       = "\033[31m"
	ColorYellow    = "\033[33m"
	ColorPurple    = "\033[35m"
	ColorLightBlue = "\033[36m"
	ColorPink      = "\033[38;5;198m"
	ColorDarkBlue  = "\033[34m"
	ColorReset     = "\033[0m"
)

func SetGlobalWidths(siteRegionWidth, rowIndexWidth int) {
	globalMaxSiteRegionWidth, globalMaxRowIndexWidth = siteRegionWidth, rowIndexWidth
}

func Setup(site, region string, rowIndex int) *Logger {
	combined := site
	if region != "" {
		combined = site + " " + region
	}

	return &Logger{
		combined: combined,
		rowIndex: rowIndex,
	}
}

func formatTime() string {
	return time.Now().Format("15:04:05.00")
}

func InitializeLogger() {
	initMu.Lock()
	defer initMu.Unlock()

	if initialised {
		return
	}

	logChan = make(chan string, channelSize)

	for i := range numWorkers {
		wgWorkers.Add(1)
		go func(workerID int) {
			defer wgWorkers.Done()
			for msg := range logChan {
				fmt.Print(msg)
			}
		}(i)
	}

	initialised = true
	shutdown = false
}

func (l *Logger) Success(format string, args ...any) {
	l.logMessage(ColorGreen, format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	l.logMessage(ColorRed, format, args...)
}

func (l *Logger) Info(format string, args ...any) {
	l.logMessage(ColorYellow, format, args...)
}

func (l *Logger) Normal(format string, args ...any) {
	l.logMessage(ColorPurple, format, args...)
}

func (l *Logger) Monitor(format string, args ...any) {
	l.logMessage(ColorLightBlue, format, args...)
}

func (l *Logger) Important(format string, args ...any) {
	l.logMessage(ColorPink, format, args...)
}

func (l *Logger) Captcha(format string, args ...any) {
	l.logMessage(ColorDarkBlue, format, args...)
}

func (l *Logger) logMessage(color, format string, args ...any) {
	timestamp := formatTime()
	shutdownMu.Lock()
	if shutdown {
		shutdownMu.Unlock()
		return
	}
	wgProducers.Add(1)
	shutdownMu.Unlock()

	message := fmt.Sprintf(format, args...)

	logEntry := fmt.Sprintf("%s[%-*s  %s  %*d] - %s%s\n",
		color,
		globalMaxSiteRegionWidth, l.combined,
		timestamp,
		globalMaxRowIndexWidth, l.rowIndex,
		message,
		ColorReset)

	logChan <- logEntry
	wgProducers.Done()
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
	initialised = false
}
