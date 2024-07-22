package apihandler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func SegmentHandler(w http.ResponseWriter, r *http.Request) {
	segment := r.URL.Query().Get("segment")
	if segment == "" {
		http.Error(w, "Missing segment", http.StatusBadRequest)
		return
	}

	// Construct the full path to the segment file
	segmentPath := filepath.Join("output/output_segs/test.mp4_720p", segment)

	// Open the segment file
	file, err := os.Open(segmentPath)
	if err != nil {
		http.Error(w, "Error opening segment file", http.StatusInternalServerError)
		log.Printf("Error opening segment file: %v", err)
		return
	}
	defer file.Close()

	// Set the content type for MPEG-2 Transport Stream
	w.Header().Set("Content-Type", "video/MP2T")

	// Copy the file content to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("Error writing segment to response: %v", err)
	}
}

func GetPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	playlistPath := filepath.Join("output/output_segs/test.mp4_720p", "playlist.m3u8")
	file, err := os.Open(playlistPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening playlist file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending playlist file: %v", err), http.StatusInternalServerError)
		return
	}
}
