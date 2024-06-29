package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"

	h "github.com/kociumba/Kopycat/handlers"
)

type SyncConfig struct {
	Interval time.Duration `json:"interval"`
	Targets  []Target      `json:"targets"`
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
	// configPath, _ := GetRelativePath()
	return &SyncConfig{}
}

func (c *SyncConfig) AddToSync(PathOrigin, PathDestination string) {
	log.Info("Adding to sync", "origin", PathOrigin, "destination", PathDestination)

	ServerConfig.Targets = append(ServerConfig.Targets, Target{
		PathOrigin:      PathOrigin,
		PathDestination: PathDestination,
	})

	err := c.SaveConfig()
	if err != nil {
		h.Clog.Error(err)
	}
}

func (c *SyncConfig) RemoveFromSync(PathOrigin, PathDestination string) {
	log.Info("Removing from sync", "origin", PathOrigin, "destination", PathDestination)

	for i, target := range ServerConfig.Targets {
		if target.PathOrigin == PathOrigin && target.PathDestination == PathDestination {
			ServerConfig.Targets = append(ServerConfig.Targets[:i], ServerConfig.Targets[i+1:]...)
			break
		}
	}

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

	file, err := os.Open(configPath)
	if err != nil {
		h.Clog.Warn("Failed to read config file", "error", err)
		if attemptCreateConfig {
			h.Clog.Info("Failed to read config, check logs and relaunch the application")
		} else {
			*c = SyncConfig{}
			c.CreateConfig(configPath, configDir)
			ServerConfig = *c
		}
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&ServerConfig)
	if err != nil {
		h.Clog.Error(err)
	}

	*c = ServerConfig

	h.Clog.Info(ServerConfig)
}

func (c *SyncConfig) ReturnTargets() []Target {
	return c.Targets
}

func (c *SyncConfig) ReturnInterval() time.Duration {
	return c.Interval
}

func (c *SyncConfig) CreateConfig(configPath, configDir string) {
	h.Clog.Info("Creating config in", "path", configPath)
	attemptCreateConfig = true

	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		h.Clog.Error(err)
		return
	}

	file, err := os.Create(configPath)
	if err != nil {
		h.Clog.Error(err)
		return
	}
	defer file.Close()

	ServerConfig = SyncConfig{
		Interval: 60 * time.Second,
		Targets:  []Target{},
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

	configPath, _ := GetRelativePath()

	h.Clog.Info("Saving config")
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(&ServerConfig)
	if err != nil {
		return err
	}

	h.Clog.Info("Config saved")

	return nil
}