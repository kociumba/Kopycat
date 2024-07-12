package mainloop

import (
	"log"
	"os"

	"github.com/kociumba/kopycat/config"
	"github.com/kociumba/kopycat/internal"
	l "github.com/kociumba/kopycat/logger"
	"github.com/kociumba/kopycat/scheduler"
)

var (
	logFile *os.File
	Clog    *log.Logger
	counter int
)

// Lesson learned: decouple this kind of shit from the main package no matter what
var S = scheduler.NewScheduler(func() {
	CheckDirs()
})

func CheckDirs() {
	// TODO: the actuall dir syncing logic
	if counter == 0 {
		// Clog.Print("\n\n")
		l.Clog.Info("Service started with", "pid", os.Getpid())

		// Make sure all the hashes are set up and paths are valid
		internal.InitialRun(config.ServerConfig)
	}

	if counter%25 == 0 {
		l.Clog.Info("Check scheduled", "log", logFile.Name(), "call", counter)
	}

	//TODO: check sync folders for changes
	// get sha256 of files or folders and comapre to last time
	// if different then sync

	counter++
}
