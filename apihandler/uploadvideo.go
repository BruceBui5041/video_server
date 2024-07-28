package apihandler

import (
	"fmt"
	"net/http"
	"strings"
	"video_server/appconst"
	"video_server/component"
	"video_server/logger"
	"video_server/messagemodel"
	"video_server/watermill"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const maxUploadSize = 1000 << 20 // 1000 MB

func UploadVideoHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set max size for the entire request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

		// Parse multipart form
		if err := c.Request.ParseMultipartForm(maxUploadSize); err != nil {
			logger.AppLogger.Error("Failed to parse multipart form", zap.Error(err))
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File too large. Max size is 1000MB"})
			return
		}

		// Get file from request
		file, header, err := c.Request.FormFile("video")
		if err != nil {
			logger.AppLogger.Error("Failed to get file from request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		// Check if the file is a video
		contentType := header.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "video/") {
			logger.AppLogger.Error("File is not a video", zap.String("contentType", contentType))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only video files are allowed"})
			return
		}

		// Extract additional fields
		videoId := c.PostForm("videoId")
		title := c.PostForm("title")
		description := c.PostForm("description")
		slug := c.PostForm("slug")

		if videoId == "" {
			videoId = uuid.New().String()
		}

		videoInfo := &messagemodel.VideoInfo{
			VideoID:     videoId,
			Title:       title,
			Description: description,
			Slug:        slug,
			S3Key:       fmt.Sprintf("%s/%s", appconst.RawVideoS3Key, slug),
		}

		err = watermill.PublishVideoUploadedEvent(appCtx, videoInfo, file)
		if err != nil {
			logger.AppLogger.Error("publish video uploaded event", zap.Error(err), zap.String("filename", slug))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Video uploaded successfully: %s", slug)})
	}
}
