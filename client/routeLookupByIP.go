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
	var amount int
	var err error
	var result []string

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
	for {
		res, e := stream.Recv()
		if e == io.EOF {
			break
		} else if e != nil {
			err = e
			break
		}
		result = append(result, res.Route)
	}
	go func(){
		defer close(out)
		for _,v := range result {
			out <- v
		}
	}()

	return amount, out, err
}