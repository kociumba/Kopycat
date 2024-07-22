package logger

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
)

var (
	executable string

	// Always assume the pointer is not at the start of the file
	//
	// Before using call:
	//
	//	MutexLog.Seek(0, io.SeekStart)
	MutexLog *logFileMutex
	LogPath  string
	LogFile  *os.File
	Clog     *log.Logger
	err      error
)

type logFileMutex struct {
	file  *os.File
	mutex *sync.Mutex
}

// Setup sets up the logger and log file relative to the executable
func Setup() *log.Logger {
	executable, err = os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	LogPath = filepath.Join(filepath.Dir(executable), "Kopycat.log")
	MutexLog, err = GetLogFile()
	if err != nil {
		log.Fatal(err)
	}

	if MutexLog == nil {
		log.Info("MutexLog is nil")
	}

	// Clean old log files to avoid cluttering the disk with useless logs
	if err = CleanOldLogs(MutexLog); err != nil {
		log.Warn(err)
	}

	// This go ommited somewhere in the last commit
	log.SetReportCaller(true)
	log.SetReportTimestamp(true)
	log.SetTimeFormat("2006-01-02 15:04:05")

	// Initialize the logger with the file output
	Clog = log.New(MutexLog)
	Clog.SetReportTimestamp(true)
	Clog.SetTimeFormat("2006-01-02 15:04:05")
	Clog.SetReportCaller(true)

	log.Info("Logging to", "path", LogPath)

	return Clog
}

// GetLogFile returns an io.Writer with a mutex on it so that different functions can access the same file
func GetLogFile() (*logFileMutex, error) {
	execDir := filepath.Dir(executable)
	logDir := filepath.Join(execDir, "logs")
	LogPath = filepath.Join(logDir, "Kopycat.log")

	if err = os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	LogFile, err = os.OpenFile(LogPath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return &logFileMutex{file: LogFile, mutex: &sync.Mutex{}}, nil
}

func (l *logFileMutex) Write(p []byte) (int, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.file.Write(p)
}

func (l *logFileMutex) Read(p []byte) (int, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.file.Read(p)
}

func (l *logFileMutex) Truncate(size int64) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.file.Truncate(size)
}

func (l *logFileMutex) Seek(offset int64, whence int) (int64, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.file.Seek(offset, whence)
}

func (l *logFileMutex) WriteString(s string) (int, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.file.WriteString(s)
}

func (l *logFileMutex) ReadIntoBuffer() (bytes.Buffer, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(l.file); err != nil {
		log.Error("Failed to read log file", "error", err)
		return *buf, err
	}

	return *buf, nil
}
