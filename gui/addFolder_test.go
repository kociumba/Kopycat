package gui

import (
	"path/filepath"
	"testing"
)

func Test_mirrorStructure(t *testing.T) {
	type args struct {
		origin            string
		destinationVolume string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test with windows separators",
			args: args{
				origin:            `C:\Users\user\gabagool`,
				destinationVolume: `D:\`,
			},
			want: filepath.Clean(`D:/Users/user/gabagool`),
		},
		{
			name: "Test with not fucked up separators",
			args: args{
				origin:            "C:/Users/user/gabagool",
				destinationVolume: "D:/",
			},
			want: filepath.Clean(`D:/Users/user/gabagool`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mirrorStructure(tt.args.origin, tt.args.destinationVolume); got != tt.want {
				t.Errorf("mirrorStructure() = %v, want %v", got, tt.want)
			}
		})
	}
}
