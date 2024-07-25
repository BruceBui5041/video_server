package main

import (
	"log"
	"net"
	"net/http"
	"video_server/apihandler"
	"video_server/grpcserver"
	"video_server/watermill"

	// You'll need to create this package
	pb "video_server/proto/video_service/video_service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	go watermill.StartSubscribers()

	// Start HTTP server
	go startHTTPServer()

	// Start gRPC server
	startGRPCServer()
}

func startHTTPServer() {
	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})
	handler := c.Handler(r)

	r.HandleFunc("/segment/playlist/{name}", apihandler.GetPlaylistHandler).Methods("GET")
	r.HandleFunc("/segment/playlist/{name}/{resolution}/{playlistName}", apihandler.GetPlaylistHandler).Methods("GET")
	r.HandleFunc("/segment", apihandler.SegmentHandler).Methods("GET")
	r.HandleFunc("/upload", apihandler.UploadVideoHandler).Methods("POST", "OPTIONS")

	srv := &http.Server{
		Handler: handler,
		Addr:    ":3000",
	}

	log.Println("Starting HTTP server on :3000")
	log.Fatal(srv.ListenAndServe())
}

func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Register your gRPC services here
	// For example:
	pb.RegisterVideoServiceServer(s, &grpcserver.VideoServiceServer{})

	log.Println("Starting gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
