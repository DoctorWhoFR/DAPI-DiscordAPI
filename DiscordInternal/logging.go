package DiscordInternal

import (
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
)

var logLevel int

const (
	loglevelTrace = 1 + iota
	loglevelDebug
	loglevelLog
	loglevelInfo
	loglevelError
)

type logMessage struct {
	Type    string        `json:"type"`
	Message interface{}   `json:"message"`
	Args    []interface{} `json:"args"`
}

func GetEnvLogLevel() {
	if logLevel == 0 {
		_logLevel, err := strconv.Atoi(os.Getenv("DEBUGGING"))

		if err != nil {
			log.Println(err)
		}

		log.Println("LOG LEVEL", _logLevel)
		logLevel = _logLevel
	}

}

// LogTrace
// DONE big refactoring needed for logging
// need more readable logging, and less
// ```
// test
// ```
// - ${line}
// - ${fullPath}
// - [x] test
// - [ ] test
// - [ ] test
// - [x] test
// - [x] test
// <!-- epic:"debugging" #tag @1.0.1 order:0 completed:2023-06-22T18:46:48.383Z -->
func LogTrace(args ...any) {
	if logLevel <= loglevelTrace {
		log.Println(color.New(color.FgCyan).Add(color.Underline).Sprintf("[%s]", "TRACE"), args)
	}
}

func LogLog(args ...any) {
	if logLevel <= loglevelLog {
		log.Println(color.New(color.FgBlue).Add(color.Underline).Sprintf("[%s]", "LOG"), args)
	}
}

func LogDebug(args ...any) {
	if logLevel <= loglevelDebug {
		log.Println(color.New(color.FgYellow).Add(color.Underline).Sprintf("[%s]", "DEBUG"), args)

	}
}

func LogInfo(args ...any) {
	if logLevel <= loglevelInfo {
		log.Println(color.New(color.FgMagenta).Add(color.Underline).Sprintf("[%s]", "INFO"), args)

	}
}

func LogError(args ...any) {
	if logLevel <= loglevelError {
		log.Println(color.New(color.FgRed).Add(color.Underline).Sprintf("[%s]", "ERROR"), args)

	}
}
