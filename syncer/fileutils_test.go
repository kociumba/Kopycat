package syncer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	usrHomeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	testCases := []struct {
		name    string
		src     string
		dst     string
		wantErr bool
	}{
		{
			name:    "Copy directory from user home directory to user home directory",
			src:     filepath.Join(usrHomeDir, "testfolder"),
			dst:     filepath.Join(usrHomeDir, "testfolder_copy"),
			wantErr: false,
		},
		{
			name:    "Copy file from user home directory to non-existing directory in user home directory",
			src:     filepath.Join(usrHomeDir, "testfolder"),
			dst:     filepath.Join(usrHomeDir, "testfolder_backup/testfolder_copy"),
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// create source
			if _, err := os.Stat(tt.src); os.IsNotExist(err) {
				if err := os.MkdirAll(tt.src, os.ModePerm); err != nil {
					t.Fatalf("Failed to create origin folder: %v", err)
				}
				defer os.RemoveAll(tt.src)
			}

			if err := Copy(tt.src, tt.dst); (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_copyFile(t *testing.T) {
	usrHomeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	tests := []struct {
		name    string
		src     string
		dst     string
		wantErr bool
	}{
		{
			name:    "Copy file from user home directory to user home directory",
			src:     filepath.Join(usrHomeDir, "testfile.txt"),
			dst:     filepath.Join(usrHomeDir, "testfile_copy"),
			wantErr: false,
		},
		{
			name:    "Copy file from user home directory to non-existing directory in user home directory",
			src:     filepath.Join(usrHomeDir, "testfile.txt"),
			dst:     filepath.Join(usrHomeDir, "testfolder_backup/testfile_copy"), // wants a folder to copy into not a file
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := os.Stat(tt.src); os.IsNotExist(err) {
				_, err = os.Create(tt.src)
				if err != nil {
					t.Fatalf("Failed to create origin file: %v", err)
				}
			}

			if err := copyFile(tt.src, tt.dst); (err != nil) != tt.wantErr {
				t.Errorf("copyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
