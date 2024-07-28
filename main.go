package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"video_server/apihandler"
	"video_server/component"
	"video_server/component/grpcserver"
	"video_server/logger"
	"video_server/model/user/usertransport"
	"video_server/watermill"

	pb "video_server/proto/video_service/video_service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecretKey := os.Getenv("JWTSecretKey")

	client, conn := grpcserver.ConnectToVideoProcessingServer()
	defer conn.Close()

	appContext := component.NewAppContext(
		connectToDB(),
		watermill.NewPubsubPublisher(),
		client,
		jwtSecretKey,
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

	logger.AppLogger.Info(
		"Successfully connected to the database",
		zap.Any("dbUser", dbUser),
		zap.Any("dbHost", dbHost),
		zap.Any("dbPort", dbPort),
		zap.Any("dbName", dbName),
	)

	return db
}

func startHTTPServer(appCtx component.AppContext) {
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 300
	r.Use(cors.New(config))

	// Define your routes
	r.GET("/segment/playlist/:name", apihandler.GetPlaylistHandler(appCtx))
	r.GET("/segment/playlist/:name/:resolution/:playlistName", apihandler.GetPlaylistHandler(appCtx))
	r.GET("/segment", apihandler.SegmentHandler(appCtx))
	r.POST("/upload", apihandler.UploadVideoHandler(appCtx))

	r.POST("/login", usertransport.Login(appCtx))
	r.POST("/register", usertransport.RegisterHandler(appCtx))

	log.Println("Starting HTTP server on :3000")
	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Register your gRPC services here
	pb.RegisterVideoProcessingServiceServer(s, &grpcserver.VideoServiceServer{})

	log.Println("Starting gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
