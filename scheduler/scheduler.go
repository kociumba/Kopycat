package scheduler

import (
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type Scheduler struct {
	intervalChan chan time.Duration
	wg           sync.WaitGroup
	callback     func()
}

// Call to create a scheduler with the callback function to use
func NewScheduler(callback func()) *Scheduler {
	return &Scheduler{
		intervalChan: make(chan time.Duration),
		callback:     callback,
	}
}

// Call to start the scheduler
func (s *Scheduler) Start() {
	s.wg.Add(1)
	go s.scheduleCheck()
}

// Call to stop the scheduler
func (s *Scheduler) Stop() {
	close(s.intervalChan)
	s.wg.Wait()
}

// Call with a time.Duration to change the interval at witch the scheduler runs
func (s *Scheduler) ChangeInterval(newInterval time.Duration) {
	s.intervalChan <- newInterval
}

func (s *Scheduler) scheduleCheck() {
	defer s.wg.Done()

	interval := time.Second * 1
	timer := time.NewTimer(interval)

	for {
		select {
		case <-timer.C:
			if s.callback != nil {
				// Calls the user callback
				s.callback()
			}
			timer.Reset(interval)

		case newInterval, ok := <-s.intervalChan:
			if !ok {
				log.Info("Stopping scheduler")
				return
			}
			log.Info("Changing interval to", "interval", newInterval)
			interval = newInterval
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(interval)
		}
	}
}
