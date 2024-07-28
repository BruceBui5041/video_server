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

func SegmentHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		videoName := c.Query("name")
		resolution := c.Query("resolution")
		segmentNumber := c.Query("number")

		if videoName == "" || resolution == "" || segmentNumber == "" {
			logger.AppLogger.Error("Missing required parameters")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
			return
		}

		key := filepath.Join("segments", videoName, resolution, fmt.Sprintf("segment_%s.ts", segmentNumber))
		vidSegment, err := storagehandler.GetFileFromCloudFrontOrS3(appconst.AWSVideoS3BuckerName, key)
		if err != nil {
			logger.AppLogger.Error("Error getting segment file", zap.Error(err), zap.String("key", key))
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting segment file: %v", err)})
			return
		}
		defer vidSegment.Close()

		c.Header("Content-Type", "video/MP2T")

		// Use c.Stream to handle the streaming of the video segment
		c.Stream(func(w io.Writer) bool {
			_, err := io.Copy(w, vidSegment)
			if err != nil {
				logger.AppLogger.Error("Error streaming segment file", zap.Error(err))
				// In case of error, we can't modify headers or status code here,
				// but we can log the error and return false to stop streaming
				return false
			}
			return false // Return false to indicate we're done streaming
		})
	}
}
