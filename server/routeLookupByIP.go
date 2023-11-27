package main

import (
	"log"
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/metadata"
	"strconv"
)

var routesAmountHeader = "routes-amount"

func(s *ServerService) RouteLookupByIP(in *pb.RouteLookupByIPRequest, stream pb.RouteParser_RouteLookupByIPServer) error {

	rt, ok := s.sessionLookup(in.Session)
	if !ok {
		log.Printf("Session lookup failure for session %s\n", in.Session)
		return status.Errorf(codes.NotFound, "No routing table file exists for session %s", in.Session)
	}
	amount, out, err := rt.FindRoutes(in.Ip, in.AllRoutes)
	if err != nil {
		return status.Errorf(codes.Internal, "Failure during route lookup for %s", in.Ip)
	}
	md := metadata.New(map[string]string{
		routesAmountHeader: strconv.Itoa(amount),
	})
	err = stream.SendHeader(md)
	if err != nil {
		return status.Error(codes.Internal, "error during sending header")
	}

	for r := range out {
		stream.Send(&pb.RouteLookupByIPResponse{Route: r.String()})
	}
	return nil
}
