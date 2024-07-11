package main

import (
	"flag"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/kardianos/service"
	"github.com/kociumba/kopycat/config"
	"github.com/kociumba/kopycat/controller"
	"github.com/kociumba/kopycat/gui"
	"github.com/kociumba/kopycat/handlers"
	"github.com/kociumba/kopycat/internal"
)

var logger service.Logger

type program struct {
	exit chan struct{}
}

var guiServer *gui.GUIServer

var port = flag.String("port", "", "Port to start the server.")

func (p *program) Start(s service.Service) error {
	// Check if running in terminal
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
	handlers.Setup()

	// Load config
	configManager := config.NewSyncConfig()
	configManager.ReadConfig()

	// Start main scheduler
	internal.S.Start()
	// Always update the interval
	if config.ServerConfig.Interval < time.Second*10 {
		config.ServerConfig.Interval = time.Second * 10
	}
	internal.S.ChangeInterval(config.ServerConfig.Interval)

	// Do not call this first or logs will get fucked
	if *port == "" {
		*port = "42069"
	}

	// Start web GUI
	guiServer = gui.NewGUIServer(*port)
	guiServer.Start()

	// start log cleaner
	// to make this actually work i would have to mutex the log file
	// and intercept the charmbracelet/log package
	// handlers.LogCleaner.Start()
}

func (p *program) Stop(service service.Service) error {
	// Stop should not block. Return with a few seconds.

	// Stop all running tasks
	internal.S.Stop()
	service.Stop()

	// Stop web GUI
	err := guiServer.Stop()
	if err != nil {
		logger.Error(err)
	}

	// Stop the service
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
