package main

import (
	"log"
	"net/http"
	"video_server/apihandler"
	"video_server/watermill"

	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	watermill.StartSubscribers()

	http.HandleFunc("/segment", enableCORS(apihandler.SegmentHandler))
	http.HandleFunc("/segment/playlist", enableCORS(apihandler.GetPlaylistHandler))
	http.HandleFunc("/upload", enableCORS(apihandler.UploadVideoHandler))

	log.Println("Starting server on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
