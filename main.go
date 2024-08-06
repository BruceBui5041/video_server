package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"video_server/apihandler"
	"video_server/common"
	"video_server/component"
	"video_server/component/cache"
	"video_server/component/grpcserver"
	"video_server/logger"
	"video_server/middleware"
	"video_server/model/category/categorytransport"
	"video_server/model/course/coursetransport"
	"video_server/model/user/usertransport"
	"video_server/model/video/videotransport"
	"video_server/watermill"

	pb "video_server/proto/video_service/video_service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	logger.CreateAppLogger()

	jwtSecretKey := os.Getenv("JWTSecretKey")

	client, conn := grpcserver.ConnectToVideoProcessingServer()
	defer conn.Close()

	awsSession, err := createAWSSession()
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	appContext := component.NewAppContext(
		connectToDB(),
		watermill.NewPubsubPublisher(),
		client,
		jwtSecretKey,
		createDynamoDBClient(awsSession),
		awsSession,
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

	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormlogger.Info, // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,           // Don't include params in the SQL log
			Colorful:                  true,            // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

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
	// config.AllowAllOrigins = true
	config.AllowOrigins = []string{"http://localhost:8080"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.MaxAge = 300
	r.Use(cors.New(config))

	// Define your routes

	video := r.Group("/video", middleware.RequiredAuth(appCtx))
	{
		// for get master list
		video.GET("/playlist/:course_slug/:video_id", apihandler.GetPlaylistHandler(appCtx))

		// for get video playlish
		video.GET(
			"/playlist/:course_slug/:video_id/:resolution/:playlistName",
			apihandler.GetPlaylistHandler(appCtx),
		)

		video.GET("", apihandler.SegmentHandler(appCtx))
		// video.GET("/:id", videotransport.GetVideoBySlug(appCtx))
		video.PUT("/:id", videotransport.UpdateVideoHandler(appCtx))
	}

	previewVideo := r.Group("/preview")
	{
		// for get master list
		previewVideo.GET("/playlist/:course_slug/:video_id", apihandler.GetPlaylistHandler(appCtx))

		// for get video playlish
		previewVideo.GET(
			"/playlist/:course_slug/:video_id/:resolution/:playlistName",
			apihandler.GetPlaylistHandler(appCtx),
		)

		previewVideo.GET("", apihandler.SegmentHandler(appCtx))
	}

	courseGroup := r.Group("/course")
	{
		courseGroup.POST("", middleware.RequiredAuth(appCtx), coursetransport.CreateCourseHandler(appCtx))
		courseGroup.GET("", coursetransport.ListCourses(appCtx))
		courseGroup.PUT("/courses/:id", coursetransport.UpdateCourseHandler(appCtx))
	}

	videoGroup := r.Group("/video")
	{
		videoGroup.POST("", middleware.RequiredAuth(appCtx), videotransport.CreateVideoHandler(appCtx))
		videoGroup.GET("/:course_slug", middleware.RequiredAuth(appCtx), videotransport.ListCourseVideos(appCtx))
	}

	categoryGroup := r.Group("/category")
	{
		categoryGroup.GET("", categorytransport.ListCategories(appCtx))
	}

	r.POST("/upload",
		middleware.RequiredAuth(appCtx),
		apihandler.UploadVideoHandler(appCtx),
	)

	r.GET("/checkauth", middleware.RequiredAuth(appCtx), func(c *gin.Context) {
		c.JSON(http.StatusOK, common.SimpleSuccessResponse("ok"))
	})

	r.POST("/login", usertransport.Login(appCtx))
	r.POST("/register", usertransport.Register(appCtx))

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

func createAWSSession() (*session.Session, error) {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	creds := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	return sess, nil
}

func createDynamoDBClient(awsSess *session.Session) *cache.DynamoDBClient {
	client, err := cache.NewDynamoDBClient(awsSess)
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}
	return client
}
