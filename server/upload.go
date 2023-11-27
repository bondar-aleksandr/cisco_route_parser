package main

import (
	"io"
	"log"
	// "os"

	"github.com/bondar-aleksandr/cisco_route_parser/parser"
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerService) FileTransfer(stream pb.RouteParser_FileTransferServer) error {
	file := NewFile()
	defer file.Close()

	for {
		req, err := stream.Recv()
		if !file.Created {
			file.SetFile(req.Platform)
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
	return stream.SendAndClose(&pb.FileUploadResponse{FileName: file.Name})
}

func (s *ServerService) parse(f *File) error {
	InfoLogger.Printf("Parsing routes for file %s", f.Name)
	r, err := f.Open()
	defer f.Close()
	if err != nil {
		return err
	}
	tableSource := parser.NewTableSource(f.Platform, r)
	allRoutes := tableSource.Parse()
	InfoLogger.Printf("Parsing routes done, found %d routes, %d unique nexthops", allRoutes.RoutesCount(), allRoutes.NHCount())
	s.newSession(f.Name, allRoutes)
	return nil
}
