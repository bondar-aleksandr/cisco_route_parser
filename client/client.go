package main

import (
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const serverAddr = "localhost:50051"
const chunkSize = uint32(1400)

type ClientService struct {
	conn *grpc.ClientConn
	client pb.RouteParserClient
	token string
	addr string
	chunkSize uint32
}

func NewClientService() *ClientService {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial %s, : %v\n", serverAddr, err)
	}
	c := pb.NewRouteParserClient(conn)
	
	return &ClientService{
		client: c,
		addr: serverAddr,
		chunkSize: chunkSize,
	}
}

func(c *ClientService) Close() error {
	return c.conn.Close()
}

func(c *ClientService) setToken(token string) {
	c.token = token
}