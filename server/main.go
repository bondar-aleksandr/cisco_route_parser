package main

import (
	"log"
	"net"
	"os"

	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
	"google.golang.org/grpc"
)

var (
	InfoLogger  *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

const addr = "0.0.0.0:50051"

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Failed to start listener!", err)
	}
	log.Printf("Listening on %s...", addr)

	s := grpc.NewServer()
	pb.RegisterRouteParserServer(s, NewServerService())
	
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}