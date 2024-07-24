package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/kociumba/kopycat/config"
	"github.com/kociumba/kopycat/gui"
	logSetup "github.com/kociumba/kopycat/logger"
	"github.com/kociumba/kopycat/scheduler"
	"github.com/kociumba/kopycat/tasks"
)

func Test_main(t *testing.T) {
	temp := os.TempDir()

	tests := []struct {
		name        string
		origin      string
		destination string
	}{
		{
			name:        "test forced first sync",
			origin:      filepath.Join(temp, "testDir"),
			destination: filepath.Join(temp, "testDir_backup"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// main refuses to work in tests so i'm going to have to start everything manually
			// go main()

			// // fucky way of seeing if the server starts
			// for i := 0; i < 10; i++ {
			// 	time.Sleep(500 * time.Millisecond)
			// 	resp, err := http.Get("http://localhost:42069")
			// 	if err == nil {
			// 		resp.Body.Close()
			// 		break
			// 	}
			// 	if i == 9 {
			// 		t.Fatal("server did not start")
			// 	}
			// }

			log.Info("paths", "origin/destination", struct {
				Origin      string `json:"origin"`
				Destination string `json:"destination"`
			}{tt.origin, tt.destination})

			// Hope this works
			mockStart(t)
			time.Sleep(2 * time.Second)

			// make sure origin is there couse this will fail if origin does not exist
			err := os.MkdirAll(tt.origin, os.ModePerm)
			if err != nil {
				t.Fatal(err)
			}

			req := gui.AddFolderRequest{
				Origin:      tt.origin,
				Destination: tt.destination,
			}
			jsonReq, err := json.Marshal(req)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := http.Post("http://localhost:42069/add-folder", "application/json", bytes.NewReader(jsonReq))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			log.Info(resp)

			_, err = os.Stat(tt.destination)
			if err != nil {
				t.Fatal(err)
			}

			mockStop(t)

			log.Info("test done")

			log.Error(os.RemoveAll(tt.origin))
			log.Error(os.RemoveAll(tt.destination))
		})
	}
}

func mockStart(t *testing.T) {
	log.Infof("mock start")
	// Do work here

	//Always call first to init the file logger
	logSetup.Setup()

	// This is unfortunetly unusable right now
	//
	// Clean old log files to avoid cluttering the disk with useless
	// Set up a scheduler to clean old log files
	LogCleaner = scheduler.NewScheduler(func() {
		log.Info("Cleaner scheduled", "log", logSetup.LogFile.Name())
		if err := logSetup.CleanOldLogs(logSetup.MutexLog); err != nil {
			logSetup.Clog.Warn(err)
		}
	})
	LogCleaner.ChangeInterval(time.Minute * 10)

	// Load config
	configManager := config.NewSyncConfig()
	configManager.ReadConfig()

	// Start main scheduler
	tasks.S.Start()
	// This get's funny if the interval is too low 💀
	if config.ServerConfig.Interval < time.Second*10 {
		config.ServerConfig.Interval = time.Second * 10
	}
	tasks.S.ChangeInterval(config.ServerConfig.Interval)

	// Do not call this first or logs will get fucked
	port := "42069"

	// Start web GUI
	guiServer = gui.NewGUIServer(port)
	guiServer.Start()

	// start log cleaner
	// to make this actually work i would have to mutex the log file
	// and intercept the charmbracelet/log package
	LogCleaner.Start()

	// Make sure the sync targets are copied when first added
	for _, target := range config.ServerConfig.Targets {
		tasks.InitialCopy(target)
	}

	t.Log("start done")
}

func mockStop(t *testing.T) {
	// Stop should not block. Return with a few seconds.
	log.Info("mock stop")

	// Stop all running tasks
	tasks.S.Stop()
	LogCleaner.Stop()

	// Save the config to file
	config.ServerConfig.SaveConfig()

	// Stop web GUI
	err := guiServer.Stop()
	if err != nil {
		log.Error(err)
	}

	t.Log("stop done")
}
