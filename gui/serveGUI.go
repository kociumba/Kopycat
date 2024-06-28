package gui

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/kociumba/Kopycat/config"
	h "github.com/kociumba/Kopycat/handlers"
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
	FolderPath string `json:"folderPath"`
	Drive      string `json:"drive"`
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
		h.Clog.Error("GUIServer is nil")
		return fmt.Errorf("GUIServer is nil")
	}

	s.mux = http.NewServeMux()
	if s.mux == nil {
		h.Clog.Error("http.NewServeMux() returned nil")
		return fmt.Errorf("http.NewServeMux() returned nil")
	}

	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.mux,
	}

	indexHTML, err := guiFiles.ReadFile("webGUI/dashboard.html")
	if err != nil {
		h.Clog.Error("error reading index.html: %v", err)
		return fmt.Errorf("error reading index.html: %v", err)
	}

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHTML)
	})

	s.mux.HandleFunc("/add-folder", s.handleAddFolder)
	s.mux.HandleFunc("/get-system-drives", s.returnSystemDrives)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if s.server == nil {
			log.Error("GUI server not initialized")
			h.Clog.Error("GUI server not initialized")
			return
		}

		log.Infof("GUI live on http://localhost:%s", s.port)
		h.Clog.Info("GUI live", "at", s.server.Addr)

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Error starting GUI server: %v", err)
			h.Clog.Errorf("Error starting GUI server: %v", err)
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
	h.Clog.Info("GUI server stopped gracefully.", "at", s.server.Addr)
	return nil
}

func (s *GUIServer) handleAddFolder(w http.ResponseWriter, r *http.Request) {
	var req AddFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.Clog.Info("Received add folder request", "drive", req.Drive, "path", req.FolderPath)

	// Clean the input paths
	req.FolderPath = filepath.Clean(req.FolderPath)
	req.Drive = filepath.Clean(req.Drive)

	// Construct the full path
	fullPath := filepath.Join(req.Drive, req.FolderPath)

	h.Clog.Info("Constructed full path", "fullPath", fullPath)

	// Check if the full path exists and is a directory
	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		h.Clog.Error("Folder not found", "path", fullPath)
		http.Error(w, "Folder not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.Clog.Error("Error stating folder path", "error", err)
		http.Error(w, "Error stating folder path", http.StatusInternalServerError)
		return
	}
	if !info.IsDir() {
		h.Clog.Error("Path is not a folder", "path", fullPath)
		http.Error(w, "Path is not a folder", http.StatusBadRequest)
		return
	}

	// Add the folder to sync
	config.AddToSync(fullPath)

	// Prepare and send the response
	res := FolderPathResponse{FullPath: fullPath}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.Clog.Error("Error encoding response", "error", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (s *GUIServer) returnSystemDrives(w http.ResponseWriter, r *http.Request) {
	drives, err := h.GetSystemDrives()
	if err != nil {
		h.Clog.Error("Error getting system drives", "error", err)
		http.Error(w, "Error getting system drives", http.StatusInternalServerError)
		return
	}

	tmpl := `
		{{range .Drives}}
		<div class="drive">
			<input type="radio" name="base-dir" value="{{.}}"><span>{{.}}</span>
		</div>
		{{end}}
	`
	t, err := template.New("drives").Parse(tmpl)
	if err != nil {
		h.Clog.Error("Error parsing template", "error", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	data := SystemDriveResponse{Drives: drives}
	var sb strings.Builder
	err = t.Execute(&sb, data)
	if err != nil {
		h.Clog.Error("Error executing template", "error", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(sb.String()))
}
