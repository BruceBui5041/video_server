package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type SegmentControl struct {
	mu       sync.Mutex
	segments map[string][]string
	access   map[string]int
}

func NewSegmentControl() *SegmentControl {
	return &SegmentControl{
		segments: make(map[string][]string),
		access:   make(map[string]int),
	}
}

func (sc *SegmentControl) AddSegment(clientID, segment string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.segments[clientID] = append(sc.segments[clientID], segment)
}

func (sc *SegmentControl) GetNextSegment(clientID string) (string, bool) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if len(sc.segments[clientID]) == 0 {
		return "", false
	}

	segment := sc.segments[clientID][0]
	sc.segments[clientID] = sc.segments[clientID][1:]
	sc.access[clientID]++

	return segment, true
}

func (sc *SegmentControl) GetAllSegments(clientID string) []string {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.segments[clientID]
}

// CORS middleware
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	}
}

func main() {
	sc := NewSegmentControl()

	// Simulate adding segments for clients
	// go func() {
	for i := range 3 {
		sc.AddSegment("client1", fmt.Sprintf("segment_%03d.ts", i))
		// time.Sleep(5 * time.Second)
	}
	// }()

	http.HandleFunc("/segment", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.URL.Query().Get("client_id")
		if clientID == "" {
			http.Error(w, "Missing client_id", http.StatusBadRequest)
			return
		}

		segment, ok := sc.GetNextSegment(clientID)
		if !ok {
			http.Error(w, "No segments available", http.StatusNotFound)
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
	}))

	http.HandleFunc("/segment/playlist.m3u8", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.URL.Query().Get("client_id")
		if clientID == "" {
			http.Error(w, "Missing client_id", http.StatusBadRequest)
			return
		}

		playlistPath := filepath.Join("output/output_segs/test.mp4_720p", "playlist.m3u8")
		file, err := os.Open(playlistPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error opening playlist file: %v", err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		w.Header().Set("Access-Control-Allow-Origin", "*") // Enable CORS

		_, err = io.Copy(w, file)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error sending playlist file: %v", err), http.StatusInternalServerError)
			return
		}
	}))

	http.HandleFunc("/stats", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		sc.mu.Lock()
		defer sc.mu.Unlock()

		fmt.Fprintf(w, "Access stats:\n")
		for clientID, count := range sc.access {
			fmt.Fprintf(w, "Client %s: %d segments served\n", clientID, count)
		}
	}))

	log.Println("Starting server on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
