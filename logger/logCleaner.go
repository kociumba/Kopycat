package logger

import (
	"bufio"
	"bytes"
	"io"
	"time"

	"github.com/charmbracelet/log"
)

const (
	logRetentionPeriod = 1 * time.Hour // Left the wrong one in last commit xd
	// logRetentionPeriod = 20 * time.Second
	dateTimeFormat = "2006-01-02 15:04:05"
)

// CleanOldLogs cleans old log entries from the log file
func CleanOldLogs(logFile *logFileMutex) error {
	// Seek to the beginning of the log file
	_, err = logFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Error("Failed to seek to the beginning of the log file", "error", err)
		return err
	}

	now := time.Now()
	cutoffTime := now.Add(-logRetentionPeriod)

	scanner := bufio.NewScanner(logFile)
	var buffer bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()

		// Skip lines that are shorter than the dateTimeFormat length
		if len(line) < len(dateTimeFormat) {
			_, err := buffer.WriteString(line + "\n")
			if err != nil {
				log.Error("Failed to write to buffer", "error", err)
				return err
			}
			continue
		}

		// Extract the timestamp from the log line
		logTimeStr := line[:len(dateTimeFormat)]
		logTime, err := time.Parse(dateTimeFormat, logTimeStr)
		if err != nil {
			log.Error("Failed to parse log timestamp", "timestamp", logTimeStr, "error", err)
			return err
		}

		formattedCutoffTime := cutoffTime.Format(dateTimeFormat)
		parsedCutoffTime, err := time.Parse(dateTimeFormat, formattedCutoffTime)
		if err != nil {
			log.Error("Failed to parse cutoff time", "cutoffTime", formattedCutoffTime, "error", err)
			return err
		}

		// Write lines with timestamps after the cutoff time
		if logTime.After(parsedCutoffTime) {
			_, err := buffer.WriteString(line + "\n")
			if err != nil {
				log.Error("Failed to write to buffer", "error", err)
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error("Failed to scan log file", "error", err)
		return err
	}

	// Clear the original log file
	err = logFile.Truncate(0)
	if err != nil {
		log.Error("Failed to truncate log file", "error", err)
		return err
	}

	_, err = logFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Error("Failed to seek to the beginning of the log file", "error", err)
		return err
	}

	// Write the buffer contents to the log file
	_, err = logFile.Write(buffer.Bytes())
	if err != nil {
		log.Error("Failed to write buffer to log file", "error", err)
		return err
	}

	return nil
}
