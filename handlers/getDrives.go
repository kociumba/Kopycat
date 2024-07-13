package handlers

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// This is kinda redundant now
//
// TODO: remove or find use
func GetSystemDrives() ([]string, error) {
	drives := []string{}
	switch runtime.GOOS {
	case "windows":
		driveLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		for _, drive := range driveLetters {
			path := string(drive) + ":\\"
			_, err := os.Stat(path)
			if err == nil {
				drives = append(drives, path)
			}
		}
	case "darwin":
		path := "/Volumes"
		files, err := os.ReadDir(path)
		if err != nil {
			return drives, err
		}
		for _, file := range files {
			if file.IsDir() {
				drives = append(drives, filepath.Join(path, file.Name()))
			}
		}
	case "linux":
		paths := []string{"/mnt", "/media"}
		for _, path := range paths {
			files, err := os.ReadDir(path)
			if err != nil && !os.IsNotExist(err) {
				return drives, err
			}
			for _, file := range files {
				if file.IsDir() {
					drives = append(drives, filepath.Join(path, file.Name()))
				}
			}
		}
	default:
		return nil, errors.New("unsupported operating system")
	}
	return drives, nil
}
