package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/abusomani/jsonhandlers"
	"github.com/charmbracelet/log"

	h "github.com/kociumba/Kopycat/handlers"
)

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
// =-  N E E D   D I F F E R E N T   C O N F I G  -=
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

type SyncConfig struct {
	Interval time.Duration `json:"interval"`
	Targets  []struct {
		Target Target `json:"target"`
	} `json:"targets"`
	cfg jsonhandlers.JSONHandler
}

type Target struct {
	PathOrigin      string `json:"path"`
	PathDestination string `json:"drive"`
}

var (
	attemptCreateConfig = false
	ServerConfig        SyncConfig
)

func NewSyncConfig() *SyncConfig {
	configPath, _ := GetRelativePath()
	cfg := jsonhandlers.New(jsonhandlers.WithFileHandler(configPath))
	return &SyncConfig{
		cfg: *cfg.SetOptions(jsonhandlers.WithFileHandler(configPath)),
	}
}

func (c *SyncConfig) AddToSync(PathOrigin, PathDestination string) {
	log.Info("Adding to sync", "origin", PathOrigin, "destination", PathDestination)

	ServerConfig.Targets = append(ServerConfig.Targets, struct {
		Target Target `json:"target"`
	}{
		Target: Target{
			PathOrigin:      PathOrigin,
			PathDestination: PathDestination,
		},
	})

	err := c.SaveConfig()
	if err != nil {
		h.Clog.Error(err)
	}
}

func GetRelativePath() (string, string) {
	executable, err := os.Executable()
	if err != nil {
		h.Clog.Fatal(err)
	}

	execDir := filepath.Dir(executable)
	configDir := filepath.Join(execDir, "config")
	configPath := filepath.Join(configDir, "config.json")

	return configPath, configDir
}

func (c *SyncConfig) ReadConfig() {
	configPath, configDir := GetRelativePath()
	log.Info("Reading config at", "path", configPath)

	err := c.cfg.Unmarshal(&ServerConfig)
	if err != nil {
		h.Clog.Warn("Failed to read config file", "error", err)
		if attemptCreateConfig {
			h.Clog.Info("Failed to read config, check logs and relaunch the application")
		} else {
			c.CreateConfig(configPath, configDir)
		}
		return
	}

	h.Clog.Info(ServerConfig)
}

func (c *SyncConfig) CreateConfig(configPath, configDir string) {
	h.Clog.Info("Creating config in", "path", configPath)
	attemptCreateConfig = true

	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		h.Clog.Error(err)
		return
	}

	f, err := os.Create(configPath)
	if err != nil {
		h.Clog.Error(err)
		return
	}
	f.Close()

	ServerConfig = SyncConfig{
		Interval: 60 * time.Second,
		Targets: []struct {
			Target Target `json:"target"`
		}{},
	}

	err = c.SaveConfig()
	if err != nil {
		h.Clog.Error("Failed to save config", err)
	}

	c.ReadConfig()
}

func (c *SyncConfig) SaveConfig() error {
	if c == nil {
		return fmt.Errorf("SyncConfig is nil")
	}

	h.Clog.Info("Saving config")
	err := c.cfg.Marshal(&ServerConfig)
	if err != nil {
		return err
	}

	h.Clog.Info("Config saved")

	return nil
}
