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

const maxUploadSize = 100 << 20 // 100 MB

func UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > maxUploadSize {
		logger.AppLogger.Error("File too large")
		http.Error(w, "File too large. Max size is 100MB", http.StatusRequestEntityTooLarge)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		logger.AppLogger.Error("Failed to parse multipart form", zap.Error(err))
		http.Error(w, "File too large. Max size is 100MB", http.StatusRequestEntityTooLarge)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		logger.AppLogger.Error("Failed to get file from request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := filepath.Base(header.Filename)
	s3Key := fmt.Sprintf("%s/%s", appconst.RawVideoS3Key, filename)

	err = storagehandler.UploadFileToS3(file, appconst.AWSVideoS3BuckerName, s3Key)
	if err != nil {
		logger.AppLogger.Error("Failed to upload file to S3", zap.Error(err), zap.String("filename", filename))
		http.Error(w, "Failed to upload video", http.StatusInternalServerError)
		return
	}

	logger.AppLogger.Info(
		"Video uploaded successfully to S3",
		zap.String("filename", filename),
		zap.String("s3Key", s3Key),
	)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Video uploaded successfully: %s", filename)
}
