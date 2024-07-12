package internal

import (
	"github.com/charmbracelet/log"
	"github.com/kociumba/kopycat/config"
	"github.com/kociumba/kopycat/syncer"
)

var (
	err error
)

func InitialRun(config config.SyncConfig) {
	for _, target := range config.Targets {
		if target.Hash == "" {
			target.Hash, err = syncer.GetHashFromPath(target.PathOrigin)
			if err != nil {
				log.Error(err)
			}
			log.Info("Set hash", "target", target.PathOrigin, "hash", target.Hash)
		}
	}
}
