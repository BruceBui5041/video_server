package grpcserver

import (
	"context"
	"fmt"
	"video_server/common"
	"video_server/logger"
	"video_server/messagemodel"
	pb "video_server/proto/video_service/video_service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToVideoProcessingServer() (pb.VideoProcessingServiceClient, *grpc.ClientConn) {
	// videoProcessorAddr := os.Getenv("VIDEO_PROCESSOR_ADDR")
	videoProcessorAddr := "video-processor:50052"
	// Set up a connection to the server.
	conn, err := grpc.NewClient(videoProcessorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.AppLogger.Fatal("did not connect:", zap.Error(err))
	}

	// Create a new client
	videoServiceClient := pb.NewVideoProcessingServiceClient(conn)

	return videoServiceClient, conn
}

func ProcessNewVideoRequest(appCtx common.AppContext, videoInfo *messagemodel.VideoInfo) error {
	// Prepare the request
	req := &pb.VideoInfo{
		VideoId: videoInfo.VideoID,
		Title:   videoInfo.Title,
		S3Key:   videoInfo.S3Key,
	}

	// Call the gRPC method
	resp, err := appCtx.GetVideoProcessingClient().ProcessNewVideoRequest(context.Background(), req)
	if err != nil {
		logger.AppLogger.Error("ProcessNewVideoRequest failed", zap.Any("req", req), zap.Error(err))
		return err
	}

	if resp.Status != codes.OK.String() {
		logger.AppLogger.Error("ProcessNewVideoRequest resp err", zap.Any("req", req), zap.Any("resp", resp))
		return fmt.Errorf("ProcessNewVideoRequest resp err. resp status %s", resp.Status)
	}

	// Handle the response
	logger.AppLogger.Error("ProcessNewVideoRequest Response", zap.Any("req", req), zap.Any("resp", resp))
	return nil
}
