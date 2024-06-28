package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/ilyakaznacheev/cleanenv"

	h "github.com/kociumba/Kopycat/handlers"
)

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
// =-  N E E D   D I F F E R E N T   C O N F I G  -=
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

type SyncConfig struct {
	interval time.Duration `yaml:"interval"`
	targets  []struct {
		target Target `yaml:"target"`
	} `yaml:"targets"`
}

type Target struct {
	PathOrigin      string `yaml:"path"`
	PathDestination string `yaml:"drive"`
}

var (
	attemptCreateConfig = false
	ServerConfig        SyncConfig
)

func AddToSync(path string) {
	log.Info("Adding to sync", "path", path)

	err := cleanenv.UpdateEnv(&ServerConfig)
	if err != nil {
		h.Clog.Error(err)
	}

}

func GetRelativePath() string {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	execDir := filepath.Dir(executable)
	configDir := filepath.Join(execDir, "config")
	configPath := filepath.Join(configDir, "config.yml")

	return configPath
}

func ReadConfig() {
	var configPath = GetRelativePath()

	h.Clog.Info("Reading config at", "path", configPath)

	err := cleanenv.ReadConfig(configPath, &ServerConfig)
	if err != nil {
		h.Clog.Warn(err)
		if attemptCreateConfig {
			h.Clog.Info("Failed to read config, check logs and relunch Kopycat")
		} else {
			Createconfig(configPath)
		}
	}
}

func Createconfig(configPath string) {
	h.Clog.Info("Creating config in ", "path", configPath)
	attemptCreateConfig = true

	f, err := os.Create(configPath)
	if err != nil {
		h.Clog.Error(err)
	}

	f.Close()

	ReadConfig()
}
