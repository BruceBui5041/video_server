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

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SegmentHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Query("video_id"))
		if err != nil {
			panic(err)
		}

		videoId := uid.GetLocalID()
		courseSlug := c.Query("course_slug")
		resolution := c.Query("resolution")
		segmentNumber := c.Query("number")

		if resolution == "" || segmentNumber == "" || courseSlug == "" {
			logger.AppLogger.Error("Missing required parameters")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
			return
		}

		cacheKey := fmt.Sprintf("%s:%s:%d", appconst.VIDEO_URL_PREFIX, courseSlug, videoId)
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

			err = dynamoDBClient.Set(cacheKey, videoURL)
			if err != nil {
				logger.AppLogger.Error("Error caching URL in DynamoDB", zap.Error(err))
			}
		}

		key := filepath.Join(
			videoURL,
			resolution,
			fmt.Sprintf("segment_%s.ts", segmentNumber),
		)

		svc := s3.New(appCtx.GetAWSSession())
		vidSegment, err := storagehandler.GetFileFromCloudFrontOrS3(svc, appconst.AWSVideoS3BuckerName, key)
		if err != nil {
			logger.AppLogger.Error("Error getting segment file", zap.Error(err), zap.String("key", key))
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting segment file: %v", err)})
			return
		}
		defer vidSegment.Close()

		c.Header("Content-Type", "video/MP2T")

		c.Stream(func(w io.Writer) bool {
			_, err := io.Copy(w, vidSegment)
			if err != nil {
				logger.AppLogger.Error("Error streaming segment file", zap.Error(err))
				return false
			}
			return false
		})
	}
}
