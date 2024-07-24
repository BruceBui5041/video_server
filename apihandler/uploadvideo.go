package apihandler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"video_server/appconst"
	"video_server/logger"
	"video_server/storagehandler"

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

	// Generate a unique filename or use a specific naming convention
	filename := filepath.Base(header.Filename)
	s3Key := fmt.Sprintf("%s/%s", appconst.RawVideoS3Key, filename)

	// Upload the file directly to S3
	err = storagehandler.UploadFileToS3(file, appconst.AWSVideoS3BuckerName, s3Key)
	if err != nil {
		logger.AppLogger.Error("Failed to upload file to S3", zap.Error(err), zap.String("filename", filename))
		http.Error(w, "Failed to upload video", http.StatusInternalServerError)
		return
	}

	logger.AppLogger.Info("Video uploaded successfully to S3", zap.String("filename", filename), zap.String("s3Key", s3Key))

	// Respond to the client
	fmt.Fprintf(w, "Video uploaded successfully: %s", filename)
}
