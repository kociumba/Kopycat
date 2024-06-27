package handlers

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

var (
	executable string
	logFile    *os.File
	clog       *log.Logger
	counter    int
	err        error
)

func CheckDirs() {
	// TODO: the actuall dir syncing logic
	if counter == 0 {
		clog.Print("\n\n")
		clog.Info("Service started with", "pid", os.Getpid())
	}

	clog.Info("Check scheduled", "log", logFile.Name(), "call", counter)

	counter++
}

// Set up the logger and log file relative to the executable
func SetupCheck() {
	log.SetReportCaller(true)

	executable, err = os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	execDir := filepath.Dir(executable)
	logDir := filepath.Join(execDir, "logs")
	logPath := filepath.Join(logDir, "Kopycat.log")

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
	clog = log.New(logFile)
	clog.SetReportTimestamp(true)
	clog.SetTimeFormat("2006-01-02 15:04:05")
	clog.SetReportCaller(true)

	log.Info("Logging to", "path", logPath)
}
