package internal

import (
	"github.com/kociumba/kopycat/handlers"
	"github.com/kociumba/kopycat/scheduler"
)

// Lesson learned: decouple this kind of shit from the main package no matter what
var S = scheduler.NewScheduler(func() {
	handlers.CheckDirs()
})
