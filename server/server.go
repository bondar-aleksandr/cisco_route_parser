package main

import (
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
	"github.com/bondar-aleksandr/cisco_route_parser/parser"
)

type ServerService struct {
	pb.RouteParserServer
	sessionCache map[string] *parser.RoutingTable
	fileCache map[string] *File
}

func NewServerService () *ServerService {
	return &ServerService{
		sessionCache: make(map[string]*parser.RoutingTable),
		fileCache: make(map[string]*File),
	}
}

func(s *ServerService) newSession (id string, rt *parser.RoutingTable, f *File) {
	s.sessionCache[id] = rt
	s.fileCache[id] = f
	InfoLogger.Printf("Added session for file %s", id)
}

func(s *ServerService) sessionLookup(id string) (*parser.RoutingTable, bool) {
	rt, ok := s.sessionCache[id]
	return rt, ok
}

func(s *ServerService) deleteSession(id string) error {
	err := s.fileCache[id].Delete()
	if err != nil {
		WarnLogger.Printf("Failed to delete %v file, reason: %v\n", s.fileCache[id].OutputFile, err)
		return err
	}
	defer delete(s.sessionCache, id)
	return nil
}