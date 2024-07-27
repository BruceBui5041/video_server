package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"video_server/apihandler"
	"video_server/common"
	"video_server/grpcserver"
	"video_server/logger"
	"video_server/messagemodel"
	"video_server/watermill"

	// You'll need to create this package
	pb "video_server/proto/video_service/video_service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func connectToDB() *gorm.DB {
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.AppLogger.Fatal(err.Error())
	}
	fmt.Println("Successfully connected to the database")
	return db
}

func startHTTPServer() {
	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"*"},
		AllowCredentials:   true,
		ExposedHeaders:     []string{"Content-Length"},
		MaxAge:             300,
		OptionsPassthrough: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		Debug: false,
	})

	appContext := common.NewAppContext(connectToDB())

	// Apply CORS middleware to all routes
	r.Use(c.Handler)

	// Define your routes
	r.HandleFunc("/segment/playlist/{name}", apihandler.GetPlaylistHandler(appContext)).Methods("GET")
	r.HandleFunc(
		"/segment/playlist/{name}/{resolution}/{playlistName}",
		apihandler.GetPlaylistHandler(appContext),
	).Methods("GET")
	r.HandleFunc("/segment", apihandler.SegmentHandler(appContext)).Methods("GET")
	r.HandleFunc("/upload", apihandler.UploadVideoHandler(appContext)).Methods("POST", "OPTIONS")
	r.HandleFunc("/test", test).Methods("GET")

	// Create server
	srv := &http.Server{
		Handler: r, // Use the router directly, not wrapped in CORS handler
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

func test(w http.ResponseWriter, r *http.Request) {
	videoInfo := &messagemodel.VideoInfo{
		VideoID: uuid.New().String(),
		Title:   "filename",
		S3Key:   "s3Key",
	}

	go watermill.PublishVideoUploadedEvent(videoInfo, nil)
}
