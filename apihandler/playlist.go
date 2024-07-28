package apihandler

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"video_server/appconst"
	"video_server/component"
	"video_server/logger"
	"video_server/storagehandler"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetPlaylistHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		videoName := c.Param("name")
		resolution := c.Param("resolution")
		playlistName := c.Param("playlistName")

		if videoName == "" {
			logger.AppLogger.Error("Missing video name")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing video name"})
			return
		}

		// Construct the key for the master playlist
		key := filepath.Join("segments", videoName, "master.m3u8")

		if playlistName != "" {
			key = filepath.Join("segments", videoName, resolution, playlistName)
		}

		playlist, err := storagehandler.GetFileFromCloudFrontOrS3(appconst.AWSVideoS3BuckerName, key)
		if err != nil {
			logger.AppLogger.Error("Error getting playlist file", zap.Error(err), zap.String("key", key))
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting playlist file: %v", err)})
			return
		}
		defer playlist.Close()

		c.Header("Content-Type", "application/vnd.apple.mpegurl")

		// Use c.Stream to handle the streaming of the playlist file
		c.Stream(func(w io.Writer) bool {
			_, err := io.Copy(w, playlist)
			if err != nil {
				logger.AppLogger.Error("Error streaming playlist file", zap.Error(err))
				// In case of error, we can't modify headers or status code here,
				// but we can log the error and return false to stop streaming
				return false
			}
			return false // Return false to indicate we're done streaming
		})
	}
}
