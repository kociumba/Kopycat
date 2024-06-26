package main

import (
	"fmt"
	"time"

	"github.com/kociumba/Kopycat/scheduler"
	"github.com/kociumba/Kopycat/service"
)

func main() {

	// platform specific setup
	service.ServiceSetup()

	scheduler.StartScheduler(check)

	// Simulate changing the interval after some time
	time.Sleep(time.Second * 5)
	scheduler.ChangeInterval(time.Second * 2) // Change interval to 1 minute

	// Simulate changing the interval again
	time.Sleep(time.Second * 5)
	scheduler.ChangeInterval(time.Second * 5) // Change interval to 30 seconds

	// Let it run for a bit before stopping
	time.Sleep(time.Second * 5)
	scheduler.StopScheduler()
}

func check() {
	fmt.Println("Checking all files at", time.Now())
}
