package main

import (
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
	"github.com/bondar-aleksandr/cisco_route_parser/parser"
	"fmt"
)

type ServerService struct {
	pb.RouteParserServer
	sessionCache map[string] *parser.RoutingTable
}

func NewServerService () *ServerService {
	return &ServerService{
		sessionCache: make(map[string]*parser.RoutingTable),
	}
}

func(s *ServerService) newSession (id string, rt *parser.RoutingTable) {
	s.sessionCache[id] = rt
	fmt.Println(s.sessionCache)
}
