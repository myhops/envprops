package cmd

import (
	"fmt"
	"log/slog"
	"strings"
)

type CopyFile struct {
	From string
	To   string
}

type copyFilesValue []CopyFile

func (c *copyFilesValue) String() string {
	if c == nil || len(*c) == 0 {
		return ""
	}
	// Reconstruct the original "from:to" strings for display
	parts := make([]string, len(*c))
	for i, cf := range *c {
		parts[i] = fmt.Sprintf("%s:%s", cf.From, cf.To)
	}
	return strings.Join(parts, ",")
}

func (c *copyFilesValue) Set(value string) error {
	logger := slog.Default().With(
		"type", "copyFilesValue",
		"method", "set",
	)

	vals := strings.Split(value, ",")
	for _, vv := range vals {
		logger.Debug("called", "value", vv)
		parts := strings.SplitN(vv, ":", 2) // Split only once
		if len(parts) != 2 {
			return fmt.Errorf("invalid copyfiles format: '%s'. Expected 'from:to'", vv)
		}
		*c = append(*c, CopyFile{From: parts[0], To: parts[1]})
	}

	return nil
}

func (c *copyFilesValue) Type() string {
	return "from:to" // A descriptive type for the flag
}

func NewCopyFilesValue(cf *[]CopyFile) *copyFilesValue {
	return (*copyFilesValue)(cf)
}
