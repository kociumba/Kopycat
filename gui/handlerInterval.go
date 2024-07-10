package gui

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kociumba/kopycat/config"
	h "github.com/kociumba/kopycat/handlers"
	"github.com/kociumba/kopycat/internal"
)

type IntervalRequest struct {
	Interval time.Duration `json:"interval"`
}

func (s *GUIServer) returnCurrentInterval(w http.ResponseWriter, r *http.Request) {
	h.Clog.Info("Returning current", "interval", config.ServerConfig.Interval)

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	if err := enc.Encode(map[string]time.Duration{
		"interval": config.ServerConfig.Interval,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *GUIServer) setNewInterval(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Interval int `json:"interval"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config.ServerConfig.Interval = time.Duration(req.Interval)
	internal.S.ChangeInterval(config.ServerConfig.Interval)
	config.ServerConfig.SaveConfig()
}
