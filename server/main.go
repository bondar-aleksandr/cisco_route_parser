package main

import (
	"io"
	"log"
	"net"
	"os"
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InfoLogger  *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

const addr = "0.0.0.0:50051"

type Server struct {
	pb.RouteParserServer
}

func (s *Server) FileTransfer(stream pb.RouteParser_FileTransferServer) error {
	file := NewFile()
	defer file.Close()

	for {
		req, err := stream.Recv()
		if !file.Created {
			file.SetFile()
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
	return stream.SendAndClose(&pb.FileUploadResponse{FileName: file.Name})
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Failed to start listener!", err)
	}
	log.Printf("Listening on %s...", addr)

	s := grpc.NewServer()
	pb.RegisterRouteParserServer(s, &Server{})
	
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}