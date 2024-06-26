package main

import (
	"flag"
	"time"

	"github.com/charmbracelet/log"
	"github.com/kardianos/service"
	"github.com/kociumba/Kopycat/handlers"
	"github.com/kociumba/Kopycat/scheduler"
)

var logger service.Logger

type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

// Actuall program not the setup
func (p *program) run() {
	logger.Infof("I'm running %v.", service.Platform())
	// Do work here
	s := scheduler.NewScheduler(func() {
		handlers.CheckDirs()
	})

	s.Start()

	// Simulate changing the interval after some time
	time.Sleep(time.Second * 5)
	s.ChangeInterval(time.Second * 2) // Change interval to 1 minute

	// Simulate changing the interval again
	time.Sleep(time.Second * 5)
	s.ChangeInterval(time.Second * 5) // Change interval to 30 seconds

	// Let it run for a bit before stopping
	time.Sleep(time.Second * 5)
	s.Stop()

}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}

func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := &service.Config{
		Name:        "Kopycat",
		DisplayName: "Kopycat",
		Description: "This is a the watcher service of Kopycat.",
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"},
		Option: options,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}
