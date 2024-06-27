package main

import (
	"flag"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/kardianos/service"
	"github.com/kociumba/Kopycat/controller"
	"github.com/kociumba/Kopycat/gui"
	"github.com/kociumba/Kopycat/handlers"
	"github.com/kociumba/Kopycat/scheduler"
)

var logger service.Logger

type program struct {
	exit chan struct{}
}

var guiServer *gui.GUIServer

var s = scheduler.NewScheduler(func() {
	handlers.CheckDirs()
})

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
	logger.Infof("Running as %v.", service.Platform())
	// Do work here

	//Always call first to init the file logger
	handlers.SetupCheck()

	s.Start()
	s.ChangeInterval(time.Second * 2) // Change interval to 1 second

	//Do not call this first or logs will get fucked
	guiServer = gui.NewGUIServer("42069")

	// logger.Error(guiServer.Start())
	guiServer.Start()
}

func (p *program) Stop(service service.Service) error {
	// Stop should not block. Return with a few seconds.
	s.Stop()
	service.Stop()

	err := guiServer.Stop()
	if err != nil {
		logger.Error(err)
	}

	logger.Info("Stopping Service!")
	close(p.exit)
	return nil
}

func main() {

	svcFlag := flag.String("service", "", "Control the system service.")
	// serverflag := flag.String("port", "", "Port to start the server.")
	flag.Parse()

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	options["OnFailure"] = "restart"
	options["RunOnLoad"] = "true"
	svcConfig := &service.Config{
		Name:        "Kopycat",
		DisplayName: "Kopycat",
		Description: "This is a the watcher service of Kopycat.",
		Option:      options,
	}

	prg := &program{}
	ser, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = ser.Logger(errs)
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

	// Handle commands to the service
	if len(os.Args) > 1 {
		command := os.Args[1]

		controller.ServiceController(command, ser)

		return
	}

	if len(*svcFlag) != 0 {
		err := service.Control(ser, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	err = ser.Run()
	if err != nil {
		logger.Error(err)
	}

}
