package controller

import (
	"github.com/charmbracelet/log"
	"github.com/kardianos/service"
)

var (
	err error
)

func ServiceController(command string, ser service.Service) {
	switch command {
	case "install":
		err := ser.Install()
		if err != nil {
			log.Warn(err)
			log.Warn("reinstalling...")
			err = ser.Stop()
			if err != nil {
				log.Debug(err)
			}
			err = ser.Uninstall()
			if err != nil {
				log.Fatal("Service could not be uninstalled, try running stop or restart first.", "error", err)
			}
			err = ser.Install()
			if err != nil {
				log.Fatal("Service could not be installed, try running stop or remove first.", "error", err)
			}
		}
		err = ser.Start()
		if err != nil {
			log.Fatal("Service could not be started, try running install or remove first.", "error", err)
		}
		log.Info("Service installed and started.")
	case "remove":
		err := ser.Stop()
		if err != nil {
			log.Debug(err)
		}
		err = ser.Uninstall()
		if err != nil {
			log.Fatal("Service could not be uninstalled, try running stop or install first.", "error", err)
		}
		log.Info("Service uninstalled.")
	case "start":
		err = ser.Start()
		if err != nil {
			log.Fatal("Service could not be started, try running install or remove first.", "error", err)
		}
		log.Info("Service started.")
	case "stop":
		err = ser.Stop()
		if err != nil {
			log.Fatal("Service could not be stopped, try running install or start first.", "error", err)
		}
		log.Info("Service stopped.")
	case "restart":
		err = ser.Restart()
		if err != nil {
			log.Fatal("Service could not be restarted, try running install or start first.", "error", err)
		}
		log.Info("Service restarted.")
	}
}
