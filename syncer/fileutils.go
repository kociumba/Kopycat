package syncer

import (
	"io"
	"os"
	"path/filepath"

	"github.com/kociumba/kopycat/config"
)

// Copies the whole directory recursively
func Copy(src string, dst string) error {
	var err error
	var fds []os.DirEntry
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Copy(srcfp, dstfp); err != nil {
				return err
			}
		} else {
			// Check if file exists in destination and compare hashes
			if dstHash, srcHash, err := CompareHashes(srcfp, dstfp); err != nil {
				return err
			} else if dstHash != srcHash {
				if err = copyFile(srcfp, dstfp); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	err = out.Sync()
	return err
}

func IsTargetInDestination(target config.Target) bool {
	folderName := filepath.Base(target.PathOrigin)

	for _, t := range config.ServerConfig.Targets {
		err := filepath.WalkDir(filepath.Dir(t.PathDestination), func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() && d.Name() == folderName {
				return filepath.SkipDir
			}
			return nil
		})
		if err == nil {
			return true
		}
	}
	return false
}
