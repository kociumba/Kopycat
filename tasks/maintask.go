package tasks

import (
	"os"
	"runtime"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/kociumba/kopycat/config"
	"github.com/kociumba/kopycat/internal"
	l "github.com/kociumba/kopycat/logger"
	"github.com/kociumba/kopycat/scheduler"
	"github.com/kociumba/kopycat/syncer"
)

var (
	counter int
)

// Lesson learned: decouple this kind of shit from the main package no matter what
func scheduleDirCheck() {
	CheckDirs()
}

var S = scheduler.NewScheduler(scheduleDirCheck)

// Technically should sync
//
// TODO: see if this couses any issues
func CheckDirs() {
	var wg sync.WaitGroup

	if counter == 0 {
		// Clog.Print("\n\n")
		l.Clog.Info("Service started with", "pid", os.Getpid())

		// Make sure all the hashes are set up and paths are valid
		internal.InitialRun(config.ServerConfig)
	}

	if counter%25 == 0 {
		l.Clog.Info("Check scheduled", "log", l.LogFile.Name(), "call", counter)
	}

	for _, target := range config.ServerConfig.Targets {
		log.Info(target)

		wg.Add(1)
		go func(t config.Target) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					l.Clog.Error("Recovered in CheckDirs", "error", r)
					log.Error(r)
				}
			}()

			s := syncer.NewSyncer(t)
			s.Sync()

			s.Free()
		}(target)
	}

	wg.Wait()

	defer runtime.GC()

	counter++
}
