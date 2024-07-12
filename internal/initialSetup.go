package internal

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/kociumba/kopycat/config"
	"github.com/kociumba/kopycat/sync"
)

var (
	err error
)

// Have to redeclare them here to despaghettify this shit
type SyncConfig struct {
	Interval time.Duration `json:"interval"`
	Targets  []Target      `json:"targets"`
}

type Target struct {
	PathOrigin      string `json:"path-origin"`
	PathDestination string `json:"path-destination"`
	Hash            string `json:"hash"`
}

func InitialRun(config config.SyncConfig) {
	for _, target := range config.Targets {
		if target.Hash == "" {
			target.Hash, err = sync.GetHashFromPath(target.PathOrigin)
			if err != nil {
				log.Error(err)
			}
			log.Info("Set hash", "target", target.PathOrigin, "hash", target.Hash)
		}
	}
}
