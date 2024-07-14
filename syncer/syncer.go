package syncer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/MShekow/directory-checksum/directory_checksum"
	"github.com/charmbracelet/log"
	"github.com/spf13/afero"

	"github.com/kociumba/kopycat/config"
	l "github.com/kociumba/kopycat/logger"
)

var (
	err error
)

type Syncer struct {
	target config.Target
}

// Call to create a new syncer
//
// Should call this on every new target
func NewSyncer(target config.Target) *Syncer {
	hash, err := GetHashFromTarget(target)
	if err != nil {
		l.Clog.Error(err)
	}

	if target.Hash == "" {
		target.Hash = hash
	}
	return &Syncer{
		target: target,
	}
}

// <---------------------------------------------------->
// I don't actually know which implementation is better
// Some more speed tests are needed
// <----------------------------------------------------->

// // This implementation does not rely on any 3rd party libraries but does not seem to be as optimized
// func GetHashFromTarget(target config.Target) (string, error) {
// 	hash := sha256.New()

// 	// Walk through the path, and generate a hash from each file.
// 	err = filepath.Walk(target.PathOrigin, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		// Ignore directories
// 		if info.IsDir() {
// 			return nil
// 		}

// 		// Open the file for reading
// 		file, err := os.Open(path)
// 		if err != nil {
// 			return err
// 		}
// 		defer file.Close()

// 		// Copy the file into the hash interface
// 		if _, err := io.CopyBuffer(hash, file, make([]byte, 4096)); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	hashInBytes := hash.Sum(nil)
// 	encodedHash := hex.EncodeToString(hashInBytes)
// 	return encodedHash, nil
// }

// # This version relies on 3rd party packages and a virtual filesystem
//
// Get a general checksum of a whole directory, which is derived from the files in it
func GetHashFromTarget(target config.Target) (string, error) {
	// Initialize an in-memory filesystem
	fs := afero.NewOsFs()

	// Scan the directory at the target path
	directory, err := directory_checksum.ScanDirectory(target.PathOrigin, fs)
	if err != nil {
		return "", err
	}

	// Compute the checksum of the directory
	checksum, err := directory.ComputeDirectoryChecksums()
	if err != nil {
		return "", err
	}

	return checksum, nil
}

// The same as
//
//	GetHashFromTarget(config.Target)
//
// But doesn't require the use of config.Target
func GetHashFromPath(path string) (string, error) {
	return GetHashFromTarget(config.Target{PathOrigin: path})
}

// Get a sha256 from a single file
//
// Used in conjunction with
//
//	GetHashFromPath() and GetHashFromTarget()
//
// To check which files to copy over when syncing
func GetHashFromFile(path string) (string, error) {
	hash := sha256.New()
	h, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer h.Close()

	if _, err := io.Copy(hash, h); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	encodedHash := hex.EncodeToString(hashInBytes)
	return encodedHash, nil
}

func CompareHashes(src, dst string) (string, string, error) {
	srcHash, err := GetHashFromFile(src)
	if err != nil {
		return "", "", err
	}

	dstHash := ""
	if _, err := os.Stat(dst); err == nil {
		if dstHash, err = GetHashFromFile(dst); err != nil {
			return "", "", err
		}
	}

	return dstHash, srcHash, nil
}

// if there is a change, return true
func (s *Syncer) CheckChanges() bool {
	currentHash, err := GetHashFromTarget(s.target)
	if err != nil {
		fmt.Println("Error calculating hash:", err)
		return false
	}

	return currentHash != s.target.Hash
}

// Only actually syncs if the hash has changed,
// Always updates the hash to the new value afterwards
//
// # No copying occurs if the hash has not changed
func (s *Syncer) Sync() {
	if s.target.Hash == "" {
		s.target.Hash, err = GetHashFromTarget(s.target)
		if err != nil {
			l.Clog.Error("Error calculating hash", "error", err)
			log.Error("Error calculating hash", "error", err)
		}

		// Update it in the config
		index := -1
		for i, target := range config.ServerConfig.Targets {
			if target.PathOrigin == s.target.PathOrigin && target.PathDestination == s.target.PathDestination {
				index = i
				break
			}
		}

		if index != -1 {
			config.ServerConfig.Targets[index].Hash = s.target.Hash
			err = config.ServerConfig.SaveConfig()
			if err != nil {
				l.Clog.Error("Error saving config", "error", err)
				log.Error("Error saving config", "error", err)
			}
		}

	}

	if !IsTargetInDestination(s.target) {
		err := Copy(s.target.PathOrigin, s.target.PathDestination)
		if err != nil {
			l.Clog.Error("Error copying", "error", err)
			log.Error("Error copying", "error", err)
		}
	}

	if s.CheckChanges() {
		l.Clog.Info("Syncing", "from", s.target.PathOrigin, "to", s.target.PathDestination)
		log.Info("Syncing", "from", s.target.PathOrigin, "to", s.target.PathDestination)

		err := Copy(s.target.PathOrigin, s.target.PathDestination)
		if err != nil {
			l.Clog.Error("Error syncing", "error", err)
			log.Error("Error syncing", "error", err)
		}
	}

	s.target.Hash, err = GetHashFromTarget(s.target)
	if err != nil {
		l.Clog.Error("Error calculating hash", "error", err)
		log.Error("Error calculating hash", "error", err)
	}
}

// Free all the memory the struct takes up and return true if successful
func (s *Syncer) Free() bool {
	runtime.SetFinalizer(s, nil)
	return true
}
