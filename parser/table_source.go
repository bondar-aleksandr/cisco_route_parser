package parser

import (
	"errors"
	"io"
	"github.com/bondar-aleksandr/cisco_route_parser/parser/entities"
	"github.com/bondar-aleksandr/cisco_route_parser/parser/ios"
	"github.com/bondar-aleksandr/cisco_route_parser/parser/nxos"
)

type Parser interface {
	Parse() *entities.RoutingTable
}

// type to define Parsing source. Platform is used to specify OS family, where output is taken from.
// Source is just io.Reader
type TableSource struct {
	Platform string
	Source io.Reader
	parser Parser
}

// Constructor. Creates *tableSource object, where p specifies platform type, s specifies reader
// to read data from.
func NewTableSource(p string, s io.Reader) (*TableSource, error) {
	var (
		platform string
		parser Parser
	)
	switch {
	case p == "ios":
		platform = p
		parser = ios.NewIosParser(s)
	case p == "nxos":
		platform = p
		parser = nxos.NewNxosParser(s)
	default:
		return nil, errors.New("wrong platform value specified")
		// ErrorLogger.Fatalf("Wrong platform value specified! Exiting...")
	}
	ts := &TableSource{
		Platform: platform,
		Source: s,
		parser: parser,
	}
	return ts, nil
}

// Run parser based on 'Platform' attribute. Returns *RoutingTable object, populated
// with values from 'Source' attribute
func(ts *TableSource) Parse() *entities.RoutingTable {
	return ts.parser.Parse()
}
