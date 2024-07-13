package gui

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/kociumba/kopycat/handlers"
	l "github.com/kociumba/kopycat/logger"
)

// I have an idea for this 😎
//
// Create an option to sync to auser picked drive by mirroring the origin path on the destination drive
func (s *GUIServer) returnSystemDrives(w http.ResponseWriter, r *http.Request) {
	drives, err := handlers.GetSystemDrives()
	if err != nil {
		l.Clog.Error("Error getting system drives", "error", err)
		http.Error(w, "Error getting system drives", http.StatusInternalServerError)
		return
	}

	tmpl := `
		{{range .Drives}}
		<div class="drive">
			<input type="radio" name="drives-option" value="{{.}}"><span>{{.}}</span>
		</div>
		{{end}}
	`
	t, err := template.New("drives").Parse(tmpl)
	if err != nil {
		l.Clog.Error("Error parsing template", "error", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	data := SystemDriveResponse{Drives: drives}
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
