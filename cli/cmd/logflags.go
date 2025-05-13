package cmd

import (
	"fmt"
	"log/slog"
	"strings"
)

// logLevelValue implements the pflag.Value interface
type logLevelValue slog.Level

// String returns the value as a string
func (l *logLevelValue) String() string {
	return slog.Level(*l).String()
}

// Set sets the value. This must be a valid value for slog.Level.
func (l *logLevelValue) Set(v string) error {
	var ll slog.Level

	if err := ll.UnmarshalText([]byte(v)); err != nil {
		return fmt.Errorf("error unmarshaling log level %s: %w", v, err)
	}
	*l = logLevelValue(ll)
	return nil
}

// Type returns the type that cobra shows in the help text
func (*logLevelValue) Type() string {
	return "slog.Level"
}

// NewLogLevelValue creates a new logLevelValue with the given default value and 
// the pointer to the Loglevel.
func NewLogLevelValue(value slog.Level, l *slog.Level) *logLevelValue {
	*l = value
	res := logLevelValue(value)
	return &res
}

// logFormatValue implements the pflag.Value interface for the logformat.
// Valid values are TEXT and JSON.
type logFormatValue string

// String returns teh value as a string
func (f *logFormatValue) String() string {
	return string(*f)
}

// Set sets the value or returns an error
func (f *logFormatValue) Set(v string) error {
	switch u := strings.ToUpper(v); u {
	case "TEXT", "JSON":
		*f = logFormatValue(u)
	default:
		return fmt.Errorf("bad log format: %s, expecting TEXT or JSON", u)
	}
	return nil
}

// Type returns the type as a string
func (f *logFormatValue) Type() string {
	return "logformat"
}

// NewLogformatValue creates a logFormatValue with the given default value.
func NewLogformatValue(value string, p *string) *logFormatValue {
	// test if the format is correct.
	*p = "TEXT"
	var f = new(logFormatValue)
	err := f.Set(value)
	if err == nil {
		*p = f.String()
	}
	return f
}

