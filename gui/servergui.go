package gui

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"sync"

	"github.com/charmbracelet/log"
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

	// Experiment with the afero.NewHttpFs might be much better than just embedding
	//
	// NOPE 💀 (at least not easly)
	indexHTML, err := guiFiles.ReadFile("webGUI/index.html")
	if err != nil {
		l.Clog.Error("error reading index.html: %v", err)
		return fmt.Errorf("error reading index.html: %v", err)
	}

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHTML)
	})

	// Create a new file system rooted at webGUI/static
	staticFs, err := fs.Sub(guiFiles, "webGUI/static")
	if err != nil {
		l.Clog.Error("error creating static file system: %v", err)
		return fmt.Errorf("error creating static file system: %v", err)
	}

	// Serve the static files using http.FileServer
	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFs))))

	s.mux.HandleFunc("/add-folder", s.handleAddFolder)
	s.mux.HandleFunc("/delete-folder", s.handleDeleteFolder)
	s.mux.HandleFunc("/get-system-drives", s.returnSystemDrives)
	s.mux.HandleFunc("/get-sync-targets", s.returnSyncTargets)
	s.mux.HandleFunc("/get-current-interval", s.returnCurrentInterval)
	s.mux.HandleFunc("/set-new-interval", s.setNewInterval)
	s.mux.HandleFunc("/get-logs", s.getLogs)

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
