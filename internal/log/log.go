package log

import (
	"os"

	"github.com/charmbracelet/log"
)

type logger struct {
	stdout *log.Logger
	stderr *log.Logger
	debug  *log.Logger
}

func (l *logger) ToStdout(path string) {
	l.stdout.Printf("cd://%s\n", path)
}

func (l *logger) Error(msg any, args ...any) {
	l.stderr.Error(msg, args...)
}

func (l *logger) Warning(msg any, args ...any) {
	l.stderr.Warn(msg, args...)
}

func (l *logger) Info(msg any, args ...any) {
	l.stdout.Info(msg, args...)
}

func (l *logger) Debug(msg any, args ...any) {
	l.debug.Debug(msg, args...)
}

func (l *logger) Fatal(msg any, args ...any) {
	l.debug.Fatal(msg, args...)
}

func newLogger() *logger {
	return &logger{
		stdout: log.NewWithOptions(os.Stdout, log.Options{
			ReportCaller: false,
		}),
		stderr: log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller: false,
		}),
		debug: log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller: true,
		}),
	}
}

// Package-level singleton
var _log = newLogger()

// Package-level functions
func ToStdout(path string)         { _log.ToStdout(path) }
func Error(msg any, args ...any)   { _log.Error(msg, args...) }
func Warning(msg any, args ...any) { _log.Warning(msg, args...) }
func Info(msg any, args ...any)    { _log.Info(msg, args...) }
func Debug(msg any, args ...any)   { _log.Debug(msg, args...) }
func Fatal(msg any, args ...any)   { _log.Fatal(msg, args...) }
