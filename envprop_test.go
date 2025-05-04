package envprops

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestPropertyLineEnv(t *testing.T) {
	type args struct {
		p      Property
		getenv func(string) string
		prefix []string
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
				p: Property{
					Key:   "foo",
					Value: "bar",
				},
				getenv: func(s string) string {
					return "foo value"
				},
			},
			want: "foo=foo value",
		},
		{
			name: "simple use default",
			args: args{
				p: Property{
					Key:   "foo",
					Value: "bar",
				},
				getenv: func(s string) string {
					return ""
				},
			},
			want: "foo=bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PropertyLineEnv(tt.args.p, tt.args.getenv, tt.args.prefix...); got != tt.want {
				t.Errorf("PropertyLineEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteProperties(t *testing.T) {
	type args struct {
		props []*Property
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "simple",
			args: args{
				props: []*Property{
					{
						Key:   "foo.val",
						Value: "bar",
					},
					{
						Key:   "baz.val",
						Value: "qux",
					},
				},
			},
			wantW:   "foo.val=bar\nbaz.val=qux\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := WriteProperties(w, tt.args.props); (err != nil) != tt.wantErr {
				t.Errorf("WriteProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("WriteProperties() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestReadProperties(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []Property
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "simple",
			args: args{
				r: bytes.NewBufferString("foo.val=bar\nbaz.val=qux\n"),
			},
			want: []Property{
				{
					Key:   "foo.val",
					Value: "bar",
				},
				{
					Key:   "baz.val",
					Value: "qux",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadProperties(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadProperties() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanupLine(t *testing.T) {
	type args struct {
		s []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{
			name: "simple",
			args: args{
				s: []byte("foo=bar"),
			},
			want: []byte("foo=bar"),
		},
		{
			name: "comment",
			args: args{
				s: []byte("foo=bar # this is a comment"),
			},
			want: []byte("foo=bar"),
		},
		{
			name: "empty",
			args: args{
				s: []byte(""),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanupLine(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cleanupLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvVarName(t *testing.T) {
	type args struct {
		name   string
		prefix []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EnvVarName(tt.args.name, tt.args.prefix...); got != tt.want {
				t.Errorf("EnvVarName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProperty_EnvVarName(t *testing.T) {
	type fields struct {
		Key   string
		Value string
	}
	type args struct {
		prefix []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "simple",
			fields: fields{
				Key:   "foo.val",
				Value: "bar",
			},
			args: args{
				prefix: []string{},
			},
			want: "FOO_VAL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Property{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := p.EnvVarName(tt.args.prefix...); got != tt.want {
				t.Errorf("Property.EnvVarName() = %v, want %v", got, tt.want)
			}
		})
	}
}
