package grpcserver

import (
	"context"
	"log"
	"os"
	"video_server/logger"
	"video_server/messagemodel"
	pb "video_server/proto/video_service/video_service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var videoServiceClient pb.VideoServiceClient

func init() {
	videoProcessorAddr := os.Getenv("VIDEO_PROCESSOR_ADDR")

	// Set up a connection to the server.
	conn, err := grpc.NewClient(videoProcessorAddr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new client
	videoServiceClient = pb.NewVideoServiceClient(conn)
}

func RequestNewVideoUploaded(videoInfo *messagemodel.VideoInfo) {
	// Prepare the request
	req := &pb.VideoInfo{
		VideoId: videoInfo.VideoID,
		Title:   videoInfo.Title,
		S3Key:   videoInfo.S3Key,
	}

	// Call the gRPC method
	resp, err := videoServiceClient.ProcessNewVideoRequest(context.Background(), req)
	if err != nil {
		logger.AppLogger.Error("ProcessNewVideoRequest failed", zap.Any("req", req))
		return
	}

	if resp.Status != codes.OK.String() {
		logger.AppLogger.Error("ProcessNewVideoRequest resp err", zap.Any("req", req), zap.Any("resp", resp))
	}

	// Handle the response
	log.Printf("Response: %v", resp)
}
