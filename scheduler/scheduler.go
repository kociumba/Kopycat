package scheduler

import (
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/kociumba/kopycat/internal"
)

// Don't use directly, to use call:
//
//	s := scheduler.NewScheduler(yourCallback)
type Scheduler struct {
	intervalChan chan time.Duration
	wg           sync.WaitGroup
	callback     func()
	once         sync.Once
}

// Call to create a scheduler with the callback function to use
//
// Things like this will work in the callback
//
//	var counter int
//	func callback() {
//		counter++
//	}()
func NewScheduler(callback func()) *Scheduler {
	return &Scheduler{
		intervalChan: make(chan time.Duration),
		callback:     callback,
	}
}

// Call to start the scheduler
func (s *Scheduler) Start() {
	s.once.Do(func() {
		s.wg.Add(1)
		go s.scheduleCheck()
	})
}

// Call to stop the scheduler
//
// Technically safe to call without first calling
//
//	Sheduler.Start()
func (s *Scheduler) Stop() {
	s.once.Do(func() {
		// Initialize the channel to avoid closing a nil channel
		s.intervalChan = make(chan time.Duration)
	})
	close(s.intervalChan)
	s.wg.Wait()
}

// Call with a time.Duration to change the interval at witch the scheduler runs
//
// Had to force it to be non blocking as it would block in some scenarios
func (s *Scheduler) ChangeInterval(newInterval time.Duration) {
	go func() {
		s.intervalChan <- newInterval
	}()
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
			log.Info("Changing interval to", "interval", newInterval, "with callback", internal.GetFuncName(s.callback, internal.Both))
			interval = newInterval
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(interval)
		}
	}
}
