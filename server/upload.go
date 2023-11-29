package main

import (
	"io"
	"log"
	"github.com/bondar-aleksandr/cisco_route_parser/parser"
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/metadata"
)

func (s *ServerService) Upload(stream pb.RouteParser_UploadServer) error {
	file := NewFile()
	defer file.OutputFile.Close()
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		WarnLogger.Println("Unable to find platform info in metadata")
		return status.Error(codes.FailedPrecondition, "no platform info in metadata")
	}
	platform := md.Get("platform")[0]

	for {
		req, err := stream.Recv()
		if !file.Created {
			err = file.SetFile(platform)
			if err != nil {
				return status.Errorf(codes.Internal, "Failed to create file: %s, reason: %v\n", file.Name, err.Error())
			}
			log.Printf("Created file %s", file.Name)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "Got error while receiving data: %v\n", err)
		}
		chunk := req.GetChunk()
		if err := file.Write(chunk); err != nil {
			file.Delete()
			return status.Errorf(codes.Internal, "Failed to write data to file: %v\n", err.Error())
		}
	}
	log.Printf("Received data for file %s", file.Name)
	err := s.parse(file)
	if err != nil {
		return status.Error(codes.Internal, "Failed to open received file")
	}
	return stream.SendAndClose(&pb.FileUploadResponse{Session: file.Name})
}

func (s *ServerService) parse(f *File) error {
	InfoLogger.Printf("Parsing routes for file %s", f.Name)
	r, err := f.Open()
	if err != nil {
		return err
	}
	tableSource := parser.NewTableSource(f.Platform, r)
	allRoutes := tableSource.Parse()
	InfoLogger.Printf("Parsing routes done, found %d routes, %d unique nexthops", allRoutes.RoutesCount(), allRoutes.NHCount())
	s.newSession(f.Name, allRoutes, f)
	return nil
}
