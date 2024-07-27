package apihandler

import (
	"fmt"
	"net/http"
	"strings"
	"video_server/appconst"
	"video_server/logger"
	"video_server/messagemodel"
	"video_server/watermill"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const maxUploadSize = 1000 << 20 // 1000 MB

func UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > maxUploadSize {
		logger.AppLogger.Error("File too large")
		http.Error(w, "File too large. Max size is 1000MB", http.StatusRequestEntityTooLarge)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		logger.AppLogger.Error("Failed to parse multipart form", zap.Error(err))
		http.Error(w, "File too large. Max size is 1000MB", http.StatusRequestEntityTooLarge)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		logger.AppLogger.Error("Failed to get file from request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if the file is a video
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "video/") {
		logger.AppLogger.Error("File is not a video", zap.String("contentType", contentType))
		http.Error(w, "Only video files are allowed", http.StatusBadRequest)
		return
	}

	// Extract additional fields
	videoId := r.FormValue("videoId")
	title := r.FormValue("title")
	description := r.FormValue("description")
	slug := r.FormValue("slug")

	if videoId == "" {
		videoId = uuid.New().String()
	}

	videoInfo := &messagemodel.VideoInfo{
		VideoID:     videoId,
		Title:       title,
		Description: description,
		Slug:        slug,
		S3Key:       fmt.Sprintf("%s/%s", appconst.RawVideoS3Key, slug),
	}

	err = watermill.PublishVideoUploadedEvent(videoInfo, file)
	if err != nil {
		logger.AppLogger.Error("publish video uploaded event", zap.Error(err), zap.String("filename", slug))
		http.Error(w, "Failed to upload video", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Video uploaded successfully: %s", slug)
}
