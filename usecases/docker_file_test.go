package usecases

import (
	"reflect"
	"testing"
)

func Test_execTemplate(t *testing.T) {
	type args struct {
		img    *Image
		tplStr string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecTemplate(tt.args.img, tt.args.tplStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("execTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("execTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
