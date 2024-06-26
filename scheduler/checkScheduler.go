package scheduler

import (
	"fmt"
	"sync"
	"time"
)

var (
	intervalChan = make(chan time.Duration)
	wg           sync.WaitGroup
)

func StartScheduler(callback func()) {
	wg.Add(1)
	go scheduleCheck(callback)
}

func StopScheduler() {
	close(intervalChan)
	wg.Wait()
}

func ChangeInterval(newInterval time.Duration) {
	intervalChan <- newInterval
}

// func checkAllFiles() {
// 	fmt.Println("Checking all files at", time.Now())
// }

func scheduleCheck(callback func()) {
	defer wg.Done()

	interval := time.Second * 1
	timer := time.NewTimer(interval)

	for {
		select {
		case <-timer.C:
			// checkAllFiles()
			callback()
			timer.Reset(interval)

		case newInterval, ok := <-intervalChan:
			if !ok {
				fmt.Println("Stopping scheduler")
				return
			}
			fmt.Println("Changing interval to", newInterval)
			interval = newInterval
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(interval)
		}
	}
}
