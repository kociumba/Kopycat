package controller

import (
	"github.com/charmbracelet/log"
	"github.com/kardianos/service"
)

type ServiceController struct {
	ser service.Service
}

func NewServiceController(ser service.Service) *ServiceController {
	return &ServiceController{ser: ser}
}

func (sc *ServiceController) Install() {
	err := sc.ser.Install()
	if err != nil {
		log.Warn(err)
		log.Warn("reinstalling...")
		err = sc.ser.Stop()
		if err != nil {
			log.Debug(err)
		}
		err = sc.ser.Uninstall()
		if err != nil {
			log.Fatal("Service could not be uninstalled, try running stop or restart first.", "error", err)
		}
		err = sc.ser.Install()
		if err != nil {
			log.Fatal("Service could not be installed, try running stop or remove first.", "error", err)
		}
	}
	err = sc.ser.Start()
	if err != nil {
		log.Fatal("Service could not be started, try running install or remove first.", "error", err)
	}
	log.Info("Service installed and started.")
}

func (sc *ServiceController) Remove() {
	err := sc.ser.Stop()
	if err != nil {
		log.Debug(err)
	}
	err = sc.ser.Uninstall()
	if err != nil {
		log.Fatal("Service could not be uninstalled, try running stop or install first.", "error", err)
	}
	log.Info("Service uninstalled.")
}

func (sc *ServiceController) Start() {
	err := sc.ser.Start()
	if err != nil {
		log.Fatal("Service could not be started, try running install or remove first.", "error", err)
	}
	log.Info("Service started.")
}

func (sc *ServiceController) Stop() {
	err := sc.ser.Stop()
	if err != nil {
		log.Fatal("Service could not be stopped, try running install or start first.", "error", err)
	}
	log.Info("Service stopped.")
}

func (sc *ServiceController) Restart() {
	err := sc.ser.Restart()
	if err != nil {
		log.Fatal("Service could not be restarted, try running install or start first.", "error", err)
	}
	log.Info("Service restarted.")
}

func (sc *ServiceController) ServiceControllerSwitch(command string) {
	switch command {
	case "install":
		sc.Install()
	case "remove":
		sc.Remove()
	case "start":
		sc.Start()
	case "stop":
		sc.Stop()
	case "restart":
		sc.Restart()
	}
}
