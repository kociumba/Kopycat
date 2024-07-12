package tasks

import (
	"os"
	"path/filepath"
	"testing"

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
				t.Fatal(err)
			}

			test1 := filepath.Join(home, "kopycat-test", "test1")
			test1backup := filepath.Join(home, "kopycat-test", "backup", "test1")
			test2 := filepath.Join(home, "kopycat-test", "test2")
			test2backup := filepath.Join(home, "kopycat-test", "backup", "test2")

			// Ensure all directories are created
			if err := os.MkdirAll(test1, 0755); err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(filepath.Join(home, "kopycat-test"))

			if err := os.MkdirAll(test2, 0755); err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(filepath.Join(home, "kopycat-test"))

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

			// Verify backup directories and files
			content1, err := os.ReadFile(filepath.Join(test1backup, "file1.txt"))
			if err != nil {
				t.Fatalf("failed to read backup file1: %v", err)
			}
			if string(content1) != "test1" {
				t.Error("file1 content mismatch")
			}

			content2, err := os.ReadFile(filepath.Join(test2backup, "file2.txt"))
			if err != nil {
				t.Fatalf("failed to read backup file2: %v", err)
			}
			if string(content2) != "test2" {
				t.Error("file2 content mismatch")
			}
		})
	}
}
