package main

import (
	"context"
	"io"
	"log"
	"strconv"
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
)

func(c *ClientService) RouteLookupByIP(ctx context.Context, ip string, all bool) (int, chan string , error) {
	out := make(chan string)
	errChan := make(chan error, 1)
	// defer close(out)
	var amount int
	var err error

	stream, err := c.client.RouteLookupByIP(ctx, &pb.RouteLookupByIPRequest{Session: c.session, Ip: ip, AllRoutes: all})
	if err != nil {
		close(out)
		return amount, out, err
	}

	md, err := stream.Header()
	if err != nil {
		log.Println("Failed to get amount of routes from gRPC header")
	}
	if amounts := md.Get("routes-amount"); len(amounts) >0 {
		amount, _ = strconv.Atoi(md.Get("routes-amount")[0]) 
	}
	go func() {
		defer close(out)
		defer close(errChan)
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				return
			} else if err != nil {
				errChan <- err
				return
			}
			out <- res.Route
		}
	}()
	err = <- errChan
	return amount, out, err
}