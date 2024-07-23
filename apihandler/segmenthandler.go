package apihandler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"video_server/appconst"
	"video_server/storagehandler"
)

// getFileFromCloudFrontOrS3 checks CloudFront first, then falls back to S3 if needed
func getFileFromCloudFrontOrS3(bucket, key string) (io.ReadCloser, error) {
	// First, try to get the file from CloudFront
	cloudfrontURL := fmt.Sprintf("https://%s/%s", appconst.AWSCloudFrontDomainName, key)
	resp, err := http.Get(cloudfrontURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		return resp.Body, nil
	}

	// If CloudFront request failed, fallback to S3
	file, err := storagehandler.GetS3File(bucket, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get file from S3: %v", err)
	}

	return file, nil
}

func SegmentHandler(w http.ResponseWriter, r *http.Request) {
	segmentNumber := r.URL.Query().Get("number")
	if segmentNumber == "" {
		http.Error(w, "Missing number", http.StatusBadRequest)
		return
	}

	videoName := r.URL.Query().Get("name")
	if videoName == "" {
		http.Error(w, "Missing video name", http.StatusBadRequest)
		return
	}

	resolution := r.URL.Query().Get("resolution")
	if resolution == "" {
		http.Error(w, "Missing video resolution", http.StatusBadRequest)
		return
	}

	key := filepath.Join("segments", videoName, resolution, fmt.Sprintf("segment_%s.ts", segmentNumber))
	vidSegment, err := getFileFromCloudFrontOrS3(appconst.AWSVideoS3BuckerName, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting segment file: %v", err), http.StatusInternalServerError)
		return
	}
	defer vidSegment.Close()

	// Set the content type for MPEG-2 Transport Stream
	w.Header().Set("Content-Type", "video/MP2T")

	_, err = io.Copy(w, vidSegment)
	if err != nil {
		log.Printf("Error writing segment to response: %v", err)
	}
}

func GetPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	videoName := r.URL.Query().Get("name")
	if videoName == "" {
		http.Error(w, "Missing video name", http.StatusBadRequest)
		return
	}

	resolution := r.URL.Query().Get("resolution")
	if resolution == "" {
		http.Error(w, "Missing video resolution", http.StatusBadRequest)
		return
	}

	key := filepath.Join("segments", videoName, resolution, fmt.Sprintf("playlist_%s.m3u8", resolution))
	playlist, err := getFileFromCloudFrontOrS3(appconst.AWSVideoS3BuckerName, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting playlist file: %v", err), http.StatusInternalServerError)
		return
	}
	defer playlist.Close()

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")

	_, err = io.Copy(w, playlist)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending playlist file: %v", err), http.StatusInternalServerError)
		return
	}
}
