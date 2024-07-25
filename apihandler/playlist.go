package apihandler

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"video_server/appconst"
	"video_server/logger"
	"video_server/storagehandler"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func GetPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoName := vars["name"]
	resolution := vars["resolution"]
	playlistName := vars["playlistName"]

	if videoName == "" {
		logger.AppLogger.Error("Missing video name")
		http.Error(w, "Missing video name", http.StatusBadRequest)
		return
	}
	// Construct the key for the master playlist
	key := filepath.Join("segments", videoName, "master.m3u8")

	if playlistName != "" {
		key = filepath.Join("segments", videoName, resolution, playlistName)
	}

	playlist, err := storagehandler.GetFileFromCloudFrontOrS3(appconst.AWSVideoS3BuckerName, key)
	if err != nil {
		logger.AppLogger.Error("Error getting playlist file", zap.Error(err), zap.String("key", key))
		http.Error(w, fmt.Sprintf("Error getting playlist file: %v", err), http.StatusInternalServerError)
		return
	}
	defer playlist.Close()

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")

	_, err = io.Copy(w, playlist)
	if err != nil {
		logger.AppLogger.Error("Error sending playlist file", zap.Error(err))
	}
}
