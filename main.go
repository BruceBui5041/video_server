package main

import (
	"log"
	"net/http"
	"video_server/apihandler"
)

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
	http.HandleFunc("/segment", enableCORS(apihandler.SegmentHandler))
	http.HandleFunc("/segment/playlist.m3u8", enableCORS(apihandler.GetPlaylistHandler))

	log.Println("Starting server on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
