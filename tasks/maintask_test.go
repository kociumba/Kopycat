package tasks

import (
	"os"
	"testing"

	"path/filepath"

	"github.com/kociumba/kopycat/config"
	"github.com/kociumba/kopycat/logger"
)

func TestCheckDirs(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test with 2 targets"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Setup()

			home, err := os.UserHomeDir()
			if err != nil {
				t.Error(err)
			}

			test1 := filepath.Clean(home + "/kopycat-test/test1")
			test1backup := filepath.Clean(home + "/kopycat-test/backup/test1")
			test2 := filepath.Clean(home + "/kopycat-test/test2")
			test2backup := filepath.Clean(home + "/kopycat-test/backup/test2")

			if err := os.MkdirAll(test1, 0755); err != nil {
				t.Error(err)
			}
			defer os.RemoveAll(test1)
			if err := os.MkdirAll(test2, 0755); err != nil {
				t.Error(err)
			}
			defer os.RemoveAll(test2)

			// Assuming they get created by the task
			defer func() {
				if err := os.RemoveAll(test1backup); err != nil {
					t.Error(err)
				}
				if err := os.RemoveAll(test2backup); err != nil {
					t.Error(err)
				}
			}()

			// Create test files
			if err := os.WriteFile(filepath.Join(test1, "file1.txt"), []byte("test1"), 0644); err != nil {
				t.Error(err)
			}
			if err := os.WriteFile(filepath.Join(test2, "file2.txt"), []byte("test2"), 0644); err != nil {
				t.Error(err)
			}

			config.ServerConfig.Targets = []config.Target{
				{
					PathOrigin:      test1,
					PathDestination: test1backup,
					Hash:            "c9e57270f7d0ddae50f5bd707c7410dc587bc3c7",
				},
				{
					PathOrigin:      test2,
					PathDestination: test2backup,
					Hash:            "c9e57270f7d0ddae50f5bd707c7410dc587bc3c7",
				},
			}

			// Run the task
			CheckDirs()

			// Check if the files are the same
			content1, err := os.ReadFile(filepath.Join(test1, "file1.txt"))
			if err != nil {
				t.Error(err)
			}
			content2, err := os.ReadFile(filepath.Join(test1backup, "file1.txt"))
			if err != nil {
				t.Error(err)
			}
			if string(content1) != string(content2) {
				t.Error("Files are not the same")
			}

			content1, err = os.ReadFile(filepath.Join(test2, "file2.txt"))
			if err != nil {
				t.Error(err)
			}
			content2, err = os.ReadFile(filepath.Join(test2backup, "file2.txt"))
			if err != nil {
				t.Error(err)
			}
			if string(content1) != string(content2) {
				t.Error("Files are not the same")
			}
		})
	}
}
