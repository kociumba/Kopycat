package tasks

import (
	"github.com/kociumba/kopycat/config"
	"github.com/kociumba/kopycat/syncer"
)

func InitialCopy(target config.Target) {
	if !syncer.IsTargetInDestination(target) {
		s := syncer.NewSyncer(target)
		s.Sync()
		s.Free()
	}
}
