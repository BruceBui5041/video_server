package apihandler

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"video_server/appconst"
	"video_server/common"
	"video_server/component"
	"video_server/logger"
	"video_server/model/course/coursestore"
	"video_server/model/video/videobiz"
	"video_server/model/video/videorepo"
	"video_server/model/video/videostore"
	"video_server/storagehandler"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetPlaylistHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("video_id"))
		if err != nil {
			panic(err)
		}

		videoId := uid.GetLocalID()
		courseSlug := c.Param("course_slug")
		resolution := c.Param("resolution")
		playlistName := c.Param("playlistName")

		if courseSlug == "" {
			logger.AppLogger.Error("Missing course slug")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing course slug"})
			return
		}

		// Check DynamoDB cache first
		cacheKey := fmt.Sprintf("%s:%d", courseSlug, videoId)
		dynamoDBClient := appCtx.GetDynamoDBClient()
		cachedURL, err := dynamoDBClient.Get(cacheKey)
		if err != nil {
			logger.AppLogger.Error("Error getting cached URL from DynamoDB", zap.Error(err))
		}

		var videoURL string
		if cachedURL != "" {
			videoURL = cachedURL
		} else {
			db := appCtx.GetMainDBConnection()
			videoStore := videostore.NewSQLStore(db)
			courseStore := coursestore.NewSQLStore(db)
			repo := videorepo.NewGetVideoRepo(videoStore, courseStore)
			biz := videobiz.NewGetVideoBiz(repo)

			video, err := biz.GetVideoById(c.Request.Context(), uint32(videoId), courseSlug)
			if err != nil {
				panic(err)
			}

			videoURL = video.VideoURL

			// Cache the videoURL in DynamoDB
			err = dynamoDBClient.Set(cacheKey, videoURL)
			if err != nil {
				logger.AppLogger.Error("Error caching URL in DynamoDB", zap.Error(err))
			}
		}

		key := filepath.Join(videoURL, "master.m3u8")

		if playlistName != "" {
			key = filepath.Join(videoURL, resolution, playlistName)
		}

		playlist, err := storagehandler.GetFileFromCloudFrontOrS3(appconst.AWSVideoS3BuckerName, key)
		if err != nil {
			logger.AppLogger.Error("Error getting playlist file", zap.Error(err), zap.String("key", key))
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting playlist file: %v", err)})
			return
		}
		defer playlist.Close()

		c.Header("Content-Type", "application/vnd.apple.mpegurl")

		c.Stream(func(w io.Writer) bool {
			_, err := io.Copy(w, playlist)
			if err != nil {
				logger.AppLogger.Error("Error streaming playlist file", zap.Error(err))
				return false
			}
			return false
		})
	}
}
