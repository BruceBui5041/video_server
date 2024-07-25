package grpcserver

import (
	"context"
	"video_server/logger"
	pb "video_server/proto/video_service/video_service" // import the generated protobuf package
)

type VideoServiceServer struct {
	pb.UnimplementedVideoServiceServer
}

func (s *VideoServiceServer) ProcessNewVideoRequest(ctx context.Context, req *pb.VideoInfo) (*pb.ProcessNewVideoResponse, error) {
	logger.AppLogger.Info("ProcessNewVideoRequest")
	return &pb.ProcessNewVideoResponse{}, nil
}
