package gui

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"sync"

	"github.com/charmbracelet/log"
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
