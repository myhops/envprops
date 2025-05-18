package envprops

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/magiconair/properties"
)

type Property struct {
	Key   string
	Value string
}

var (
	ErrFormatError = fmt.Errorf("format error")
)

func EnvVarName(name string, prefix ...string) string {
	if len(prefix) > 1 {
		panic("max one prefix allowed")
	}
	if prefix != nil && len(prefix[0]) == 0 {
		prefix = nil
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

func ReadProperties(r io.Reader) ([]*Property, error) {
	p, err := properties.LoadReader(r, properties.UTF8)
	if err != nil {
		return nil, err
	}

	res := make([]*Property, 0, p.Len()+10)
	for _, k := range p.Keys() {
		v, ok := p.Get(k)
		if !ok {
			v = ""
		}
		res = append(res, &Property{
			Key:   k,
			Value: v,
		})
	}
	return res, nil
}

func WriteProperties(w io.Writer, props []*Property) error {
	p := properties.NewProperties()

	for _, pp := range props {
		if pp.Value == "" {
			continue
		}
		p.Set(pp.Key, pp.Value)
	}
	if _, err := p.Write(w, properties.UTF8); err != nil {
		return fmt.Errorf("error writen properties")
	}
	return nil
}
