package main

import (
	"context"
	"log"
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serverAddr = "localhost:50051"
const chunkSize = uint32(1400)

type ClientService struct {
	conn *grpc.ClientConn
	client pb.RouteParserClient
	session string
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

func(c *ClientService) setSession(s string) {
	c.session = s
}

func(c *ClientService) sessionClose(ctx context.Context) {
	_, err := c.client.Close(ctx, &pb.SessionCloseRequest{Session: c.session})
	if err != nil {
		WarnLogger.Println(err)
	}
}