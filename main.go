package main

import (
	"log"
	"net/http"
	"video_server/apihandler"
	"video_server/watermill"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	go watermill.StartSubscribers()
	// go redishander.StartRedisSubscribers(redishander.RedisClient)

	r := mux.NewRouter()

	r.Use(enableCORS)

	r.HandleFunc("/segment/playlist/{name}", apihandler.GetPlaylistHandler).Methods("GET")
	r.HandleFunc("/segment/playlist/{name}/{resolution}/{playlistName}", apihandler.GetPlaylistHandler).Methods("GET")
	r.HandleFunc("/segment", apihandler.SegmentHandler).Methods("GET")
	r.HandleFunc("/upload", apihandler.UploadVideoHandler).Methods("GET")

	srv := &http.Server{
		Handler: r,
		Addr:    ":3000",
	}

	log.Println("Starting server on :3000")
	log.Fatal(srv.ListenAndServe())
}
