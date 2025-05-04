package envprops

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type Property struct {
	Key   string
	Value string
}

var (
	propSep = []byte("=")

	ErrFormatError = fmt.Errorf("format error")
)

func EnvVarName(name string, prefix ...string) string {
	if len(prefix) > 1 {
		panic("max one prefix allowed")
	}
	name = strings.ToUpper(strings.ReplaceAll(name, ".", "_"))
	return strings.Join(append(prefix, name), "_")
}

func (p Property) EnvVarName(prefix ...string) string {
	return EnvVarName(p.Key, prefix...)
}

func PropertyLineEnv(p Property, getenv func(string) string, prefix ...string) string {
	// Get the env var name
	k := EnvVarName(p.Key, prefix...)
	// Get the value
	v := getenv(k)
	if v == "" {
		v = p.Value
	}
	return p.Key + "=" + v
}

func PropertyLine(p Property, prefix ...string) string {
	return PropertyLineEnv(p, os.Getenv, prefix...)
}

func (p *Property) Unmarshal(s []byte) error {
	parts := bytes.Split(s, propSep)
	if len(parts) != 2 {
		return ErrFormatError
	}
	p.Key = string(parts[0])
	p.Value = string(parts[1])
	return nil
}

func Marshal(p *Property) ([]byte, error) {
	return []byte(p.Key + "=" + p.Value), nil
}

func cleanupLine(s []byte) []byte {
	i := bytes.Index(s, []byte("#"))
	if i >= 0 {
		s = s[:i]
	}		

	s = bytes.TrimSpace(s)
	if len(s) == 0 {
		return nil
	}
	return s
}

func ReadProperties(r io.Reader) ([]*Property, error) {
	scanner := bufio.NewScanner(r)
	var props []*Property
	for scanner.Scan() {
		b := cleanupLine(scanner.Bytes())
		if b == nil {
			continue
		}
		p := &Property{}
		if err := p.Unmarshal(b); err != nil {
			return nil, err
		}
		props = append(props, p)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return props, nil
}

func WriteProperties(w io.Writer, props []*Property) error {
	for _, p := range props {
		b, err := Marshal(p)
		if err != nil {
			return err
		}
		if _, err := w.Write(b); err != nil {
			return err
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}
	return nil
}
