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

func SegmentHandler(w http.ResponseWriter, r *http.Request) {
	videoName := r.URL.Query().Get("name")
	resolution := r.URL.Query().Get("resolution")
	segmentNumber := r.URL.Query().Get("number")

	if videoName == "" || resolution == "" || segmentNumber == "" {
		logger.AppLogger.Error("Missing required parameters")
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	key := filepath.Join("segments", videoName, resolution, fmt.Sprintf("segment_%s.ts", segmentNumber))
	vidSegment, err := storagehandler.GetFileFromCloudFrontOrS3(appconst.AWSVideoS3BuckerName, key)
	if err != nil {
		logger.AppLogger.Error("Error getting segment file", zap.Error(err), zap.String("key", key))
		http.Error(w, fmt.Sprintf("Error getting segment file: %v", err), http.StatusInternalServerError)
		return
	}
	defer vidSegment.Close()

	w.Header().Set("Content-Type", "video/MP2T")

	_, err = io.Copy(w, vidSegment)
	if err != nil {
		logger.AppLogger.Error("Error writing segment to response", zap.Error(err))
	}
}
