package handlers

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/kociumba/kopycat/scheduler"
)

var (
	executable string
	logFile    *os.File
	Clog       *log.Logger
	counter    int
	err        error

	LogCleaner *scheduler.Scheduler
)

func CheckDirs() {
	// TODO: the actuall dir syncing logic
	if counter == 0 {
		// Clog.Print("\n\n")
		Clog.Info("Service started with", "pid", os.Getpid())
	}

	if counter%25 == 0 {
		Clog.Info("Check scheduled", "log", logFile.Name(), "call", counter)
	}

	//TODO: check sync folders for changes
	// get sha256 of files or folders and comapre to last time
	// if different then sync

	counter++
}

// Set up the logger and log file relative to the executable
func Setup() *log.Logger {
	log.SetReportCaller(true)

	executable, err = os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	execDir := filepath.Dir(executable)
	logDir := filepath.Join(execDir, "logs")
	logPath := filepath.Join(logDir, "Kopycat.log")

	// I should run this continuously in scheduler but for that i would need a mutex on the log file
	//
	// Clean old log files to avoid cluttering the disk with useless logs
	if err = cleanOldLogs(logPath); err != nil {
		log.Fatal(err)
	}

	if err = os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal(err)
	}

	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the logger with the file output
	Clog = log.New(logFile)
	Clog.SetReportTimestamp(true)
	Clog.SetTimeFormat("2006-01-02 15:04:05")
	Clog.SetReportCaller(true)

	// This is unfortunetly unusable right now
	//
	// Clean old log files to avoid cluttering the disk with useless
	// Set up a scheduler to clean old log files
	// LogCleaner = scheduler.NewScheduler(func() {
	// 	if err = cleanOldLogs(logPath); err != nil {
	// 		Clog.Warn(err)
	// 	}
	// })
	// LogCleaner.ChangeInterval(time.Minute * 5)

	log.Info("Logging to", "path", logPath)

	return Clog
}
