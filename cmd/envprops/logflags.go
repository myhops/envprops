package main

import (
	"fmt"
	"log/slog"
	"strings"
)

type logLevelValue slog.Level

func (l *logLevelValue) String() string {
	return slog.Level(*l).String()
}

func (l *logLevelValue) Set(v string) error {
	var ll slog.Level

	if err := ll.UnmarshalText([]byte(v)); err != nil {
		return fmt.Errorf("error unmarshaling log level %s: %w", v, err)
	}
	*l = logLevelValue(ll)
	return nil
}

func (*logLevelValue) Type() string {
	return "slog.Level"
}

func NewLogLevelValue(value slog.Level, l *slog.Level) *logLevelValue {
	*l = value
	res := logLevelValue(value)
	return &res
}


type logFormatValue string

func (f *logFormatValue) String() string {
	return string(*f)
}

func (f *logFormatValue) Set(v string) error {
	switch u := strings.ToUpper(v); u {
	case "TEXT", "JSON":
		*f = logFormatValue(u)
	default:
		return fmt.Errorf("bad log format: %s, expecting TEXT or JSON", u)
	}
	return nil
}

func (f *logFormatValue) Type() string {
	return "logformat"
}

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

