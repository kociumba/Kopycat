package gui

import (
	"html/template"
	"net/http"
	"runtime"
	"strings"

	"github.com/kociumba/kopycat/config"
	l "github.com/kociumba/kopycat/logger"
)

func (s *GUIServer) returnSyncTargets(w http.ResponseWriter, r *http.Request) {
	data := config.NewSyncConfig()
	data.ReadConfig()

	targets := data.ReturnTargets()

	tmpl := `
	{{range .}}
		<div class="target-item">
			<span>{{.PathOrigin}} -> {{.PathDestination}}
			<button class="button" onclick="deleteTarget('{{.PathOrigin}}', '{{.PathDestination}}')">Delete</button>
			</span>
		</div>
	{{end}}
	`

	t, err := template.New("sync").Parse(tmpl)
	if err != nil {
		l.Clog.Error("Error parsing template", "error", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// h.Clog.Info("Targets:", "targets", targets)

	var sb strings.Builder
	err = t.Execute(&sb, targets)
	if err != nil {
		l.Clog.Error("Error executing template", "error", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(sb.String()))

	//Testing to see if this works
	go runtime.GC()
}
