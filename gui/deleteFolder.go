package gui

import (
	"encoding/json"
	"net/http"

	"github.com/kociumba/kopycat/config"
	l "github.com/kociumba/kopycat/logger"
)

func (s *GUIServer) handleDeleteFolder(w http.ResponseWriter, r *http.Request) {
	req := AddFolderRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		l.Clog.Error("Error decoding delete request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l.Clog.Info("Received delete folder request", "origin", req.Origin, "destination", req.Destination)

	config.ServerConfig.RemoveFromSync(req.Origin, req.Destination)
}
