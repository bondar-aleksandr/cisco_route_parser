package main

import (
	"context"
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func(s *ServerService) Close(ctx context.Context, req *pb.SessionCloseRequest) (*pb.Empty, error) {
	err := s.deleteSession(req.Session)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete routing-table file, reason: %v", err)
	}
	return &pb.Empty{}, nil
}