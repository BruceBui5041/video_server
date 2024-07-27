package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"video_server/apihandler"
	"video_server/common"
	"video_server/grpcserver"
	"video_server/logger"
	"video_server/watermill"

	// You'll need to create this package
	pb "video_server/proto/video_service/video_service"

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

	client, conn := grpcserver.ConnectToVideoProcessingServer()
	defer conn.Close()

	appContext := common.NewAppContext(
		connectToDB(),
		watermill.NewPubsubPublisher(),
		client,
	)

	go watermill.StartSubscribers(appContext)

	// Start HTTP server
	go startHTTPServer(appContext)

	// Start gRPC server
	startGRPCServer()
}

func connectToDB() *gorm.DB {
	// Get database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.AppLogger.Fatal(err.Error())
	}
	fmt.Println("Successfully connected to the database")
	return db
}

func startHTTPServer(appCtx common.AppContext) {
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

	// Apply CORS middleware to all routes
	r.Use(c.Handler)

	// Define your routes
	r.HandleFunc("/segment/playlist/{name}", apihandler.GetPlaylistHandler(appCtx)).Methods("GET")
	r.HandleFunc(
		"/segment/playlist/{name}/{resolution}/{playlistName}",
		apihandler.GetPlaylistHandler(appCtx),
	).Methods("GET")
	r.HandleFunc("/segment", apihandler.SegmentHandler(appCtx)).Methods("GET")
	r.HandleFunc("/upload", apihandler.UploadVideoHandler(appCtx)).Methods("POST", "OPTIONS")

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
	pb.RegisterVideoProcessingServiceServer(s, &grpcserver.VideoServiceServer{})

	log.Println("Starting gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
