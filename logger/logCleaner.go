package logger

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/log"
)

const (
	// logRetentionPeriod = 1 * 24 * time.Hour // Retain logs for 1 day
	logRetentionPeriod = 10 * time.Minute // testing times
	dateTimeFormat     = "2006-01-02 15:04:05"
)

func cleanOldLogs(logPath string) error {
	tempPath := logPath + ".tmp"

	logFile, err := os.OpenFile(logPath, os.O_RDONLY, 0666)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// No log file exists yet, nothing to clean
			return nil
		}
		log.Error(err)
		return err
	}

	tempFile, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = logFile.Seek(0, io.SeekStart) // Move cursor to the beginning of the file
	if err != nil {
		log.Error("Failed to seek to the beginning of the log file", "error", err)
		return err
	}

	now := time.Now()
	cutoffTime := now.Add(-logRetentionPeriod)

	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip lines that are shorter than the dateTimeFormat length
		if len(line) < len(dateTimeFormat) {
			_, err := tempFile.WriteString(line + "\n")
			if err != nil {
				log.Error("Failed to write to temp log file", "error", err)
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

		// log.Info("Log timestamp", "timestamp", logTime, "Cutoff timestamp", cutoffTime)

		// Write lines with timestamps after the cutoff time
		if logTime.After(parsedCutoffTime) {
			_, err := tempFile.WriteString(line + "\n")
			if err != nil {
				log.Error("Failed to write to temp log file", "error", err)
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error("Failed to scan log file", "error", err)
		return err
	}

	tempFile.Close()

	tempFile, err = os.OpenFile(tempPath, os.O_RDWR, 0666)
	if err != nil {
		log.Error(err)
		return err
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(tempFile); err != nil {
		log.Error("Failed to read from temp log file", "error", err)
		return err
	}
	tempFileContents := buf.String()

	trimmedTempFileContents := strings.TrimLeftFunc(tempFileContents, func(r rune) bool {
		return unicode.IsSpace(r)
	})

	if _, err := tempFile.WriteString(trimmedTempFileContents); err != nil {
		log.Error("Failed to write to temp log file", "error", err)
		return err
	}

	// Close before rename
	logFile.Close()
	tempFile.Close()

	if err := os.Rename(tempPath, logPath); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
