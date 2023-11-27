package main

import (
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
	"github.com/bondar-aleksandr/cisco_route_parser/parser"
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
}

func(s *ServerService) sessionLookup(id string) (*parser.RoutingTable, bool) {
	rt, ok := s.sessionCache[id]
	return rt, ok
}