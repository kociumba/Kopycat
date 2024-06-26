//go:build !windows

package service

import (
	"github.com/charmbracelet/log"
)

// Placeholder untill other platforms are supported
func ServiceSetup() {
	log.Debug("Not running Windows setup")
}
