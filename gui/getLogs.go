package gui

import (
	"io"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/kociumba/kopycat/logger"
)

// Return the content of the log file as text/plain.
//
// If logger.MutexLog is nil, return an error.
func (s *GUIServer) getLogs(w http.ResponseWriter, r *http.Request) {
	// logger.Clog.Info("Received get logs request", "Method", r.Method)

	if logger.MutexLog == nil {
		http.Error(w, "logger.MutexLog is nil", http.StatusInternalServerError)
		return
	}

	// Seek to the beginning of the log file
	_, err := logger.MutexLog.Seek(0, io.SeekStart)
	if err != nil {
		log.Error("Failed to seek to the beginning of the log file", "error", err)
		return
	}

	logs, err := logger.MutexLog.ReadIntoBuffer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write(logs.Bytes()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
