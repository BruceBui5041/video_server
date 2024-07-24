package apihandler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"video_server/logger"

	"go.uber.org/zap"
)

func UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	// Set a reasonable maximum upload size (e.g., 100MB)
	r.ParseMultipartForm(100 << 20)

	// Get the file from the request
	file, header, err := r.FormFile("video")
	if err != nil {
		logger.AppLogger.Error("Failed to get file from request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the uploads directory if it doesn't exist
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		logger.AppLogger.Error("Failed to create uploads directory", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(filepath.Join("uploads", header.Filename))
	if err != nil {
		logger.AppLogger.Error("Failed to create destination file", zap.Error(err), zap.String("filename", header.Filename))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem
	_, err = io.Copy(dst, file)
	if err != nil {
		logger.AppLogger.Error("Failed to copy uploaded file", zap.Error(err), zap.String("filename", header.Filename))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.AppLogger.Info("Video uploaded successfully", zap.String("filename", header.Filename))

	// Respond to the client
	fmt.Fprintf(w, "Video uploaded successfully: %s", header.Filename)
}
