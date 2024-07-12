package gui

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/kociumba/kopycat/config"
	l "github.com/kociumba/kopycat/logger"
	"github.com/kociumba/kopycat/mainloop"
)

type IntervalRequest struct {
	Interval time.Duration `json:"interval"`
}

func (s *GUIServer) returnCurrentInterval(w http.ResponseWriter, r *http.Request) {
	data := map[string]time.Duration{
		"interval": config.ServerConfig.Interval,
	}
	tmpl := `
		<span id="current-interval" class="text">Sync interval: {{.interval}}</span>
	`
	t, err := template.New("sync").Parse(tmpl)
	if err != nil {
		l.Clog.Error("Error parsing template", "error", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	var sb strings.Builder
	err = t.Execute(&sb, data)
	if err != nil {
		l.Clog.Error("Error executing template", "error", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(sb.String()))
}

func (s *GUIServer) setNewInterval(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Interval int `json:"interval"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		l.Clog.Error("Error decoding request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config.ServerConfig.Interval = time.Duration(req.Interval)
	mainloop.S.ChangeInterval(config.ServerConfig.Interval)
	config.ServerConfig.SaveConfig()
}
