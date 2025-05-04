package main

import (
	"github/myhops/envprops"
	"testing"
)

func Test_getEnvVars(t *testing.T) {
	type args struct {
		props  []*envprops.Property
		getenv func(string) string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "simple",
			args: args{
				props: []*envprops.Property{
					{
						Key:   "foo.val",
						Value: "bar",
					},
				},
				getenv: func(s string) string {
					switch s {
					case "FOO_VAL":
						return "foo value"
					default:
						return ""
					}
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getEnvVars(tt.args.props, tt.args.getenv)
			t.Log("done")
		})
	}
}
