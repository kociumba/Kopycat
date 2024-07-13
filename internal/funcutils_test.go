package internal

import "testing"

func TestGetFuncName(t *testing.T) {
	type args struct {
		f   interface{}
		arg []Args
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Test with empty args", args{f: TestGetFuncName, arg: []Args{}}, "github.com/kociumba/kopycat/internal.TestGetFuncName"},
		{"Test with args", args{f: TestGetFuncName, arg: []Args{OnlyFunc}}, "TestGetFuncName"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFuncName(tt.args.f, tt.args.arg...); got != tt.want {
				t.Errorf("GetFuncName() = %v, want %v", got, tt.want)
			}
		})
	}
}
