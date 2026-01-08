package log

import (
	"os"

	"github.com/charmbracelet/log"
)

type logger struct {
	stdout *log.Logger
	stderr *log.Logger
	debug *log.Logger
}

func (l *logger) CD(path string) {
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

func New() *logger {
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