package gui

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/kociumba/kopycat/config"
	"github.com/kociumba/kopycat/handlers"
	l "github.com/kociumba/kopycat/logger"
)

//go:embed webGUI/*
var guiFiles embed.FS

type GUIServer struct {
	mux    *http.ServeMux
	server *http.Server
	wg     sync.WaitGroup
	port   string
}

type AddFolderRequest struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

type FolderPathResponse struct {
	FullPath string `json:"fullPath"`
}

type SystemDriveResponse struct {
	Drives []string `json:"drives"`
}

func NewGUIServer(port string) *GUIServer {
	return &GUIServer{
		port: port,
		// clog: handlers.SetupCheck(),
	}
}

func (s *GUIServer) Start() error {
	if s == nil {
		l.Clog.Error("GUIServer is nil")
		return fmt.Errorf("GUIServer is nil")
	}

	s.mux = http.NewServeMux()
	if s.mux == nil {
		l.Clog.Error("http.NewServeMux() returned nil")
		return fmt.Errorf("http.NewServeMux() returned nil")
	}

	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.mux,
	}

	indexHTML, err := guiFiles.ReadFile("webGUI/dashboard.html")
	if err != nil {
		l.Clog.Error("error reading index.html: %v", err)
		return fmt.Errorf("error reading index.html: %v", err)
	}

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHTML)
	})

	s.mux.HandleFunc("/add-folder", s.handleAddFolder)
	s.mux.HandleFunc("/delete-folder", s.handleDeleteFolder)
	s.mux.HandleFunc("/get-system-drives", s.returnSystemDrives)
	s.mux.HandleFunc("/get-sync-targets", s.returnSyncTargets)
	s.mux.HandleFunc("/get-current-interval", s.returnCurrentInterval)
	s.mux.HandleFunc("/set-new-interval", s.setNewInterval)

	// FOR PROFILING
	// s.mux.HandleFunc("/debug/pprof/", pprof.Index)
	// s.mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	// s.mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// s.mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// s.mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// s.mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	// s.mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	// s.mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	// s.mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if s.server == nil {
			log.Error("GUI server not initialized")
			l.Clog.Error("GUI server not initialized")
			return
		}

		log.Infof("GUI live on http://localhost:%s", s.port)
		l.Clog.Info("GUI live", "at", s.server.Addr)

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Error starting GUI server: %v", err)
			l.Clog.Errorf("Error starting GUI server: %v", err)
		}
	}()

	return nil
}

func (s *GUIServer) Stop() error {
	if s.server == nil {
		return fmt.Errorf("server not initialized")
	}

	// Shutdown the HTTP server gracefully
	err := s.server.Shutdown(context.TODO())
	if err != nil {
		return fmt.Errorf("error shutting down server: %v", err)
	}

	// Wait for all goroutines to finish
	s.wg.Wait()

	log.Info("GUI server stopped.")
	l.Clog.Info("GUI server stopped gracefully.", "at", s.server.Addr)
	return nil
}

// DEPRECATED
//
// TODO: remove
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