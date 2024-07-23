package gui

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kociumba/kopycat/config"
	"github.com/kociumba/kopycat/handlers"
	l "github.com/kociumba/kopycat/logger"
	"github.com/kociumba/kopycat/syncer"
	"github.com/kociumba/kopycat/tasks"
)

func (s *GUIServer) handleAddFolder(w http.ResponseWriter, r *http.Request) {
	var req AddFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l.Clog.Info("Received add folder request", "origin", req.Origin, "destination", req.Destination)

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
		l.Clog.Error("Origin folder not found", "path", req.Origin)
		http.Error(w, "Origin folder not found", http.StatusNotFound)
		return
	}
	if err != nil {
		l.Clog.Error("Error resolving origin folder path", "error", err)
		http.Error(w, "Error resolving origin folder path", http.StatusInternalServerError)
		return
	}
	if !info.IsDir() {
		l.Clog.Error("Origin path is not a folder", "path", req.Origin)
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
			l.Clog.Error("Error resolving folder path", "error", err)
			http.Error(w, "Error resolving folder path", http.StatusInternalServerError)
			return
		} else {
			break
		}
	}
	if checkPath == filepath.VolumeName(req.Origin) {
		l.Clog.Error("Destination folder not found", "path", req.Destination)
		http.Error(w, "Destination folder not found", http.StatusNotFound)
		return
	}

	// Propably the best way i can do this without creating a circular dependency
	hash, err := syncer.GetHashFromTarget(config.Target{
		PathOrigin:      req.Origin,
		PathDestination: req.Destination,
	})
	if err != nil {
		l.Clog.Error("Error getting hash from target", "error", err)
		http.Error(w, "Error getting hash from target", http.StatusInternalServerError)
		return
	}

	// Check if destination is a volume name
	volumes, err := handlers.GetSystemDrives()
	if err != nil {
		l.Clog.Error("Error getting system drives", "error", err)
		http.Error(w, "Error getting system drives", http.StatusInternalServerError)
		return
	}
	if runtime.GOOS == "windows" {
		for _, volume := range volumes {
			if req.Destination == volume {
				l.Clog.Info("Destination is a volume name", "path", req.Destination)
				req.Destination = mirrorStructure(req.Origin, req.Destination)
				break
			}
		}
	} else {
		l.Clog.Error("Volume mirroring is currently only supported on windows", "error", err)
		http.Error(w, "Volume mirroring is currently only supported on windows", http.StatusInternalServerError)
		return
	}

	if req.Origin == req.Destination {
		l.Clog.Error("Origin and destination are the same", "origin", req.Origin, "destination", req.Destination)
		http.Error(w, "Origin and destination are the same", http.StatusBadRequest)
		return
	}

	// Add the folder to sync
	config.ServerConfig.AddToSync(req.Origin, req.Destination, hash)
	// config.ServerConfig.SaveConfig()

	// Make sure the files are copied on addition
	tasks.InitialCopy(config.Target{
		PathOrigin:      req.Origin,
		PathDestination: req.Destination,
		Hash:            hash,
	})

	// Prepare and send the response
	res := FolderPathResponse{FullPath: req.Origin + " -> " + req.Destination}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		l.Clog.Error("Error encoding response", "error", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// returns the mirrored pathwith / on unix and \ on windows
func mirrorStructure(origin, destinationVolume string) string {

	originVolume := filepath.VolumeName(origin)

	// l.Clog.Info("the found volume", "volume", originVolume)

	if originVolume == "" {
		return origin
	}

	parts := strings.Split(filepath.ToSlash(origin), "/")
	for i, part := range parts {
		if part == originVolume {
			parts[i] = destinationVolume
			break
		}
	}

	return filepath.Clean(filepath.ToSlash((strings.Join(parts, "/"))))
}
