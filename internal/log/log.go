package log

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"

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

// Fmt formats a template string with flexible placeholder syntax:
//   - {}         → next positional argument
//   - {0}, {1}   → indexed argument
//   - {Name}     → struct field value
func Fmt(template string, args ...any) string {
	re := regexp.MustCompile(`\{(\w*)\}`)
	posIndex := 0

	return re.ReplaceAllStringFunc(template, func(match string) string {
		inner := match[1 : len(match)-1]

		// Empty {} → positional
		if inner == "" {
			if posIndex < len(args) {
				val := args[posIndex]
				posIndex++
				return fmt.Sprint(val)
			}
			return match
		}

		// Numeric {0}, {1} → indexed
		if idx, err := strconv.Atoi(inner); err == nil {
			if idx < len(args) {
				return fmt.Sprint(args[idx])
			}
			return match
		}

		// Named {FieldName} → struct field lookup
		for _, arg := range args {
			v := reflect.ValueOf(arg)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			if v.Kind() == reflect.Struct {
				field := v.FieldByName(inner)
				if field.IsValid() {
					return fmt.Sprint(field.Interface())
				}
			}
		}

		return match // not found, leave placeholder as-is
	})
}
