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
			name:    "Test with 1 target",
			src:     filepath.Join(usrHomeDir, "testFolder"),
			dst:     filepath.Join(usrHomeDir, "testFolder_backup"),
			wantErr: false,
		},
		{
			name:    "Test depth",
			src:     filepath.Join(usrHomeDir, "testfolder"),
			dst:     filepath.Join(usrHomeDir, "testfolder_backup/subfolder/subfolder2"),
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := os.MkdirAll(tt.src, 0755)
			if err != nil {
				t.Errorf("Failed to create test directory: %v", err)
				return
			}
			defer os.RemoveAll(tt.src)

			err = Copy(tt.src, tt.dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if _, err := os.Stat(tt.dst); os.IsNotExist(err) {
				t.Errorf("Copy() did not create the destination directory: %v", tt.dst)
				return
			}

			// Check if destination contains the same files or the folders match up to the source.
			err = filepath.Walk(tt.src, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				relPath, err := filepath.Rel(tt.src, path)
				if err != nil {
					return err
				}
				dstPath := filepath.Join(tt.dst, relPath)
				dstInfo, err := os.Stat(dstPath)
				if os.IsNotExist(err) {
					t.Errorf("Copy() did not copy file %s to destination %s", path, dstPath)
					return nil
				}
				if err != nil {
					return err
				}
				if info.IsDir() != dstInfo.IsDir() {
					t.Errorf("Copy() did not correctly create directory %s", dstPath)
				}
				if !info.IsDir() && info.Size() != dstInfo.Size() {
					t.Errorf("Copy() did not correctly copy file %s to destination %s", path, dstPath)
				}
				return nil
			})
			if err != nil {
				t.Errorf("Copy() failed to walk destination directory: %v", err)
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
			name:    "Test 1 file",
			src:     filepath.Join(usrHomeDir, "testfile.txt"),
			dst:     filepath.Join(usrHomeDir, "testfile_backup.txt"),
			wantErr: false,
		},
		{
			name:    "Test one file with folder",
			src:     filepath.Join(usrHomeDir, "testfile.txt"),
			dst:     filepath.Join(usrHomeDir, "testfolder_backup/testfile_copy.txt"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := os.Stat(tt.src); os.IsNotExist(err) {
				f, err := os.Create(tt.src)
				if err != nil {
					t.Fatalf("Failed to create origin file: %v", err)
				}
				defer f.Close()
				f.WriteString("Test")
				defer os.Remove(tt.src)
			}

			if err := copyFile(tt.src, tt.dst); (err != nil) != tt.wantErr {
				t.Errorf("copyFile() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				srcHash, err := GetHashFromFile(tt.src)
				if err != nil {
					t.Errorf("Failed to get hash from src file: %v", err)
				}
				dstHash, err := GetHashFromFile(tt.dst)
				if err != nil {
					t.Errorf("Failed to get hash from dst file: %v", err)
				}
				if srcHash != dstHash {
					t.Errorf("src and dst files are different")
				}
			}
		})
	}
}
