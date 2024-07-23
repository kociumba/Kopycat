package gui

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/charmbracelet/log"
)

func Test_mirrorStructure(t *testing.T) {
	type args struct {
		origin            string
		destinationVolume string
	}
	tests := []struct {
		name          string
		args          args
		want          string
		onlyOnWindows bool
		onlyOnLinux   bool
	}{
		{
			name: "test with windows separators",
			args: args{
				origin:            `C:\Users\user\gabagool`,
				destinationVolume: `D:\`,
			},
			want:          filepath.Clean(`D:/Users/user/gabagool`),
			onlyOnWindows: true,
		},
		{
			name: "Test with not fucked up separators",
			args: args{
				origin:            "C:/Users/user/gabagool",
				destinationVolume: `D:\`,
			},
			want:          filepath.Clean("D:/Users/user/gabagool"),
			onlyOnWindows: true,
		},
		{
			name: "test with linux volume d",
			args: args{
				origin:            "/home/user/gabagool",
				destinationVolume: "/mnt/d",
			},
			want:        filepath.Clean("/mnt/d/home/user/gabagool"),
			onlyOnLinux: true,
		},
		{
			name: "test with linux volume e",
			args: args{
				origin:            "/home/user/otherstuff",
				destinationVolume: "/mnt/e",
			},
			want:        filepath.Clean("/mnt/e/home/user/otherstuff"),
			onlyOnLinux: true,
		},
	}
	for _, tt := range tests {
		// Skip tests based on the OS
		if (tt.onlyOnWindows && runtime.GOOS != "windows") || (tt.onlyOnLinux && runtime.GOOS != "linux") {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {

			log.Info("found volume", "volume", filepath.VolumeName(tt.args.origin))

			if got := mirrorStructure(tt.args.origin, tt.args.destinationVolume); got != tt.want {
				t.Errorf("mirrorStructure() = %v, want %v", got, tt.want)
			}
		})
	}
}
