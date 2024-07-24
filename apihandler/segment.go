package apihandler

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"video_server/appconst"
	"video_server/logger"
	"video_server/storagehandler"

	"go.uber.org/zap"
)

// getFileFromCloudFrontOrS3 checks CloudFront first, then falls back to S3 if needed
func getFileFromCloudFrontOrS3(bucket, key string) (io.ReadCloser, error) {
	// First, try to get the file from CloudFront
	cloudfrontURL := fmt.Sprintf("http://%s/%s", appconst.AWSCloudFrontDomainName, key)
	resp, err := http.Get(cloudfrontURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		logger.AppLogger.Info("File retrieved from CloudFront", zap.String("key", key))
		return resp.Body, nil
	}

	// If CloudFront request failed, fallback to S3
	file, err := storagehandler.GetS3File(bucket, key)
	if err != nil {
		logger.AppLogger.Error("Failed to get file from S3", zap.Error(err), zap.String("bucket", bucket), zap.String("key", key))
		return nil, fmt.Errorf("failed to get file from S3: %v", err)
	}

	logger.AppLogger.Info("File retrieved from S3", zap.String("bucket", bucket), zap.String("key", key))
	return file, nil
}

func SegmentHandler(w http.ResponseWriter, r *http.Request) {
	segmentNumber := r.URL.Query().Get("number")
	if segmentNumber == "" {
		logger.AppLogger.Error("Missing segment number")
		http.Error(w, "Missing number", http.StatusBadRequest)
		return
	}

	videoName := r.URL.Query().Get("name")
	if videoName == "" {
		logger.AppLogger.Error("Missing video name")
		http.Error(w, "Missing video name", http.StatusBadRequest)
		return
	}

	resolution := r.URL.Query().Get("resolution")
	if resolution == "" {
		logger.AppLogger.Error("Missing video resolution")
		http.Error(w, "Missing video resolution", http.StatusBadRequest)
		return
	}

	key := filepath.Join("segments", videoName, resolution, fmt.Sprintf("segment_%s.ts", segmentNumber))
	vidSegment, err := getFileFromCloudFrontOrS3(appconst.AWSVideoS3BuckerName, key)
	if err != nil {
		logger.AppLogger.Error("Error getting segment file", zap.Error(err), zap.String("key", key))
		http.Error(w, fmt.Sprintf("Error getting segment file: %v", err), http.StatusInternalServerError)
		return
	}
	defer vidSegment.Close()

	// Set the content type for MPEG-2 Transport Stream
	w.Header().Set("Content-Type", "video/MP2T")

	_, err = io.Copy(w, vidSegment)
	if err != nil {
		logger.AppLogger.Error("Error writing segment to response", zap.Error(err))
	}
}

func GetPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	videoName := r.URL.Query().Get("name")
	if videoName == "" {
		logger.AppLogger.Error("Missing video name")
		http.Error(w, "Missing video name", http.StatusBadRequest)
		return
	}

	resolution := r.URL.Query().Get("resolution")
	if resolution == "" {
		logger.AppLogger.Error("Missing video resolution")
		http.Error(w, "Missing video resolution", http.StatusBadRequest)
		return
	}

	key := filepath.Join("segments", videoName, resolution, fmt.Sprintf("playlist_%s.m3u8", resolution))
	playlist, err := getFileFromCloudFrontOrS3(appconst.AWSVideoS3BuckerName, key)
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
		http.Error(w, fmt.Sprintf("Error sending playlist file: %v", err), http.StatusInternalServerError)
		return
	}
}
