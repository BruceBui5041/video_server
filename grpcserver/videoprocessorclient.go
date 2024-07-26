package grpcserver

import (
	"context"
	"video_server/logger"
	"video_server/messagemodel"
	pb "video_server/proto/video_service/video_service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

func OpenVideoProcessorClient() (pb.VideoServiceClient, *grpc.ClientConn) {
	// videoProcessorAddr := os.Getenv("VIDEO_PROCESSOR_ADDR")
	videoProcessorAddr := "video-processor:50052"
	// Set up a connection to the server.
	conn, err := grpc.NewClient(videoProcessorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.AppLogger.Fatal("did not connect:", zap.Error(err))
	}

	// Create a new client
	videoServiceClient := pb.NewVideoServiceClient(conn)

	return videoServiceClient, conn
}

func RequestNewVideoUploaded(videoInfo *messagemodel.VideoInfo) {
	videoServiceClient, conn := OpenVideoProcessorClient()
	defer conn.Close()

	// Prepare the request
	req := &pb.VideoInfo{
		VideoId: videoInfo.VideoID,
		Title:   videoInfo.Title,
		S3Key:   videoInfo.S3Key,
	}

	// Call the gRPC method
	resp, err := videoServiceClient.ProcessNewVideoRequest(context.Background(), req)
	if err != nil {
		logger.AppLogger.Error("ProcessNewVideoRequest failed", zap.Any("req", req), zap.Error(err))
		return
	}

	if resp.Status != codes.OK.String() {
		logger.AppLogger.Error("ProcessNewVideoRequest resp err", zap.Any("req", req), zap.Any("resp", resp))
	}

	// Handle the response
	logger.AppLogger.Error("ProcessNewVideoRequest Response", zap.Any("req", req), zap.Any("resp", resp))
}
