package gui

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kociumba/kopycat/config"
	h "github.com/kociumba/kopycat/handlers"
	"github.com/kociumba/kopycat/sync"
)

func (s *GUIServer) handleAddFolder(w http.ResponseWriter, r *http.Request) {
	var req AddFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.Clog.Info("Received add folder request", "origin", req.Origin, "destination", req.Destination)

	// Clean the input paths
	req.Origin = filepath.Clean(req.Origin)
	req.Destination = filepath.Clean(req.Destination)

	// // Construct the full path
	// fullPath := filepath.Join(req.Drive, req.FolderPath)

	// h.Clog.Info("Constructed full path", "fullPath", fullPath)

	// Check if the full path exists and is a directory
	// Check if the origin path exists and is a directory
	info, err := os.Stat(req.Origin)
	if os.IsNotExist(err) {
		h.Clog.Error("Origin folder not found", "path", req.Origin)
		http.Error(w, "Origin folder not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.Clog.Error("Error resolving origin folder path", "error", err)
		http.Error(w, "Error resolving origin folder path", http.StatusInternalServerError)
		return
	}
	if !info.IsDir() {
		h.Clog.Error("Origin path is not a folder", "path", req.Origin)
		http.Error(w, "Origin path is not a folder", http.StatusBadRequest)
		return
	}

	// recursively check if os stat returns an error if the error is os.IsNotExist check the next layer down and so on until root
	var checkPath string
	for checkPath != filepath.VolumeName(req.Destination) {
		checkPath = filepath.Dir(checkPath)
		_, err = os.Stat(checkPath)
		if err != nil && os.IsNotExist(err) {
			continue
		} else if err != nil {
			h.Clog.Error("Error resolving folder path", "error", err)
			http.Error(w, "Error resolving folder path", http.StatusInternalServerError)
			return
		} else {
			break
		}
	}
	if checkPath == filepath.VolumeName(req.Origin) {
		h.Clog.Error("Destination folder not found", "path", req.Destination)
		http.Error(w, "Destination folder not found", http.StatusNotFound)
		return
	}

	// Propably the best way i can do this without creating a circular dependency
	hash, err := sync.GetHashFromTarget(config.Target{
		PathOrigin:      req.Origin,
		PathDestination: req.Destination,
	})
	if err != nil {
		h.Clog.Error("Error getting hash from target", "error", err)
		http.Error(w, "Error getting hash from target", http.StatusInternalServerError)
		return
	}

	// Add the folder to sync
	config.ServerConfig.AddToSync(req.Origin, req.Destination, hash)
	// config.ServerConfig.SaveConfig()

	// Prepare and send the response
	res := FolderPathResponse{FullPath: req.Origin + " -> " + req.Destination}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.Clog.Error("Error encoding response", "error", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
