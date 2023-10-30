package parser

import (
	"fmt"
	"log"
	"net/netip"
	"os"
	"sort"
	"strings"
	"github.com/mitchellh/hashstructure/v2"
	"io"
)

var (
	infoLogger  *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// type to define Parsing source. Platform is used to specify OS family, where output is taken from.
// Source is just io.Reader
type tableSource struct {
	Platform string
	Source io.Reader
	parser func(*tableSource) *RoutingTable
}

// Constructor. Creates *tableSource object, where p specifies platform type, s specifies reader
// to read data from.
func NewTableSource(p string, s io.Reader) *tableSource {
	var platform string
	var parser func(*tableSource) *RoutingTable
	switch {
	case p == "ios":
		platform = p
		parser = parseRouteIOS
	case p == "nxos":
		platform = p
		parser = parseRouteNXOS
	default:
		errorLogger.Fatalf("Wrong platform value specified! Exiting...")
	}
	return &tableSource{
		Platform: platform,
		Source: s,
		parser: parser,
	}
}

// runs corresponding parser
func(ts *tableSource) Parse() *RoutingTable {
	return ts.parser(ts)
}

//Single route entity
type route struct {
	Network netip.Prefix
	Type string
	NHList []uint64
	ParentRT *RoutingTable
}

func (r *route) String() string {
	nhlist := []*nextHop{}
	for _,v := range r.NHList {
		nhlist = append(nhlist, r.ParentRT.NH[v])
	}
	return fmt.Sprintf("%s route to %s network via %v", r.Type, r.Network.String(), nhlist)
}

// Constructor
func newRoute(rt *RoutingTable) *route {
	return &route{
		NHList: make([]uint64, 0),
		ParentRT: rt,
	}
}

// to reset attributes to default values, except 'ParentRT'
// TODO: find out how to setup default value for r.Network
func (r *route) reset() {
	r.NHList = make([]uint64, 0)
	r.Type = ""
	// r.Network = nil
}

// whenever NH is added to the Route object, it's also added to RoutingTable object
func (r *route) addNextHop(nh *nextHop) {
	r.ParentRT.addNextHop(nh)
	r.NHList = append(r.NHList, nh.getHash())
}

// next-hops count
func (r *route) nhCount() int {
	return len(r.NHList)
}

// Nexthop entity
type nextHop struct {
	IsIP bool
	Addr netip.Addr `hash:"string"`
	Intf string
}

// Constructor
func newNextHop(s string) *nextHop {
	if v , err := netip.ParseAddr(s); err != nil {
		return &nextHop{IsIP: false, Intf: s}
	} else {
		return &nextHop{IsIP: true, Addr: v}
	}
}

func (nh *nextHop) String() string {
	if nh.IsIP {
		return fmt.Sprintf("{NextHop: %s}", nh.Addr)
	}
	return fmt.Sprintf("{NextHop: %s}", nh.Intf)
}

func (nh *nextHop) getHash() (uint64) {
	hash, err := hashstructure.Hash(nh, hashstructure.FormatV2, nil)
	if err != nil {
		errorLogger.Fatalf("Cannot compute hash from nexthop %s due to: %q", nh.String(), err)
	}
	return hash
}

// Routing table type. Consists of *Route slice and *nextHop map. Only unique next-hops are stored.
// Next-hops stored in map, where keys are their hashes, values are hext-hops themselves.
// TODO: add table name (vrf), add to parser as well
type RoutingTable struct{
	Table string
	Routes []*route
	NH map[uint64]*nextHop
}

// Constructor for Routing table
func newRoutingTable(table string) *RoutingTable {
	return &RoutingTable{
		Table: table,
		Routes: make([]*route, 0),
		NH: make(map[uint64]*nextHop),
	}
}

func (rt *RoutingTable) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("Table: %s\n", rt.Table))
	b.WriteString("Routes:\n")
	for _,v := range rt.Routes {
		b.WriteString(v.String() + "\n")
	}
	b.WriteString("Next-Hops:\n")
	for k,v := range rt.NH {
		b.WriteString(fmt.Sprintf("%d : %v\n", k, v))
	}
	return b.String()
}

func (rt *RoutingTable) addRoute(r *route) {
	r.ParentRT = rt
	rt.Routes = append(rt.Routes, r)
}

// Only unique values added
func (rt *RoutingTable) addNextHop(nh *nextHop) {
	if _, ok := rt.NH[nh.getHash()]; ok {
		return
	}
	rt.NH[nh.getHash()] = nh
}

func (rt *RoutingTable) RoutesCount() int {
	return len(rt.Routes)
}

func (rt *RoutingTable) NHCount() int {
	return len(rt.NH)
}

func (rt *RoutingTable) getLast() *route {
	return rt.Routes[rt.RoutesCount() - 1]
}

// For test purposes. It's assumed that there is only one route to the destination in routing table
func (rt *RoutingTable) getByNetwork(s string) *route {
	netw, err := netip.ParsePrefix(s)
	if err != nil {
		errorLogger.Printf("Cannot parse ip %s", s)
		return nil
	}
	for _,v := range rt.Routes {
		if netw.String() == v.Network.String() {
			return v
		}
	}
	return nil
}

// FindRoutes func return channel of *Route objects, which contain "ip" specified.
// Routes put in channel are ordered based on prefix lenght, starting from more specific
// If "all" flag is specified, func return all matched routes, otherwise only best match returned
func (rt *RoutingTable) FindRoutes(ip string, all bool) (<-chan *route, error) {
	out := make(chan *route)
	parsedIp, err := netip.ParseAddr(ip)
	if err != nil {
		close(out)
		return out, err
	}

	indexes := []*route{}
	for _, v := range rt.Routes {
		if v.Network.Contains(parsedIp) {
			indexes = append(indexes, v)
		}
	}
	// if no routes found
	if len(indexes) == 0 {
		close(out)
		return out, nil
	}
	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i].Network.Bits() > indexes[j].Network.Bits()
	})

	// for case when we need to return only best route
	if !all {
		indexes = indexes[:1]
	}
	go func(){
		defer close(out)
		for _,v := range indexes {
			out <- v
		}
	}()
	return out, nil
}

// FindRoutesByNH func finds all routes with specified nexthop.
// Result returned as a channel
func (rt *RoutingTable) FindRoutesByNH(n string) <-chan *route {
	out := make(chan *route)
	nh := newNextHop(n)
	res := []*route{}
	for _, route := range rt.Routes {
		for _, v := range route.NHList {
			if v == nh.getHash() {
				res = append(res, route)
			}
		}
	}
	// if no routes found
	if len(res) == 0 {
		close(out)
		return out
	}
	go func(){
		defer close(out)
		for _,v := range res {
			out <- v
		}
	}()
	return out
}

// FindUniqNexthop func finds all unique NextHop objects. Result returned as a channel.
// "ipOnly" flag gives ability to specify whether we need to get only NextHops with IP address
func (rt *RoutingTable) FindUniqNexthops (ipOnly bool) <-chan *nextHop {
	out := make(chan *nextHop)
	nhList := []*nextHop{}
	for _, nh := range rt.NH {
		if ipOnly {
			if nh.IsIP {
				nhList = append(nhList, nh)
			}
		} else {
			nhList = append(nhList, nh)
		}
	}
	go func(){
		defer close(out)
		for _,v := range nhList {
			out <-v
		}
	}()
	return out
}