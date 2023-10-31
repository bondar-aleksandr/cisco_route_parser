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
}

// Constructor. Creates *tableSource object, where p specifies platform type, s specifies reader
// to read data from.
func NewTableSource(p string, s io.Reader) *tableSource {
	var platform string
	switch {
	case p == "ios":
		platform = p
	case p == "nxos":
		platform = p
	default:
		errorLogger.Fatalf("Wrong platform value specified! Exiting...")
	}
	return &tableSource{
		Platform: platform,
		Source: s,
	}
}

// Run parser based on 'Platform' attribute. Returns *RoutingTable object, populated
// with values from 'Source' attribute
func(ts *tableSource) Parse() *RoutingTable {
	switch ts.Platform {
	case "ios":
		return parseRouteIOS(ts)
	case "nxos":
		return parseRouteNXOS(ts)
	}
	return nil
}

// Type to describe route entity. 'ParentRT' attribute refers to *RoutingTable
// this route belongs. NHList attribute stores next-hop hashes, which are used by
// this route.
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

// Returns pointer to new route object with 'ParentRT' attribute set up to rt value
func newRoute(rt *RoutingTable) *route {
	return &route{
		NHList: make([]uint64, 0),
		ParentRT: rt,
	}
}

// to add nextHop to route. Whenever NH is added to the route object, 
// it's also added to RoutingTable object
func (r *route) addNextHop(nh *nextHop) {
	r.ParentRT.addNextHop(nh)
	r.NHList = append(r.NHList, nh.getHash())
}

// next-hops count
func (r *route) nhCount() int {
	return len(r.NHList)
}

// Type to describe nexthop entity
type nextHop struct {
	IsIP bool
	Addr netip.Addr `hash:"string"`
	Intf string
}

// Returns pointer to new nextHop object. Attributes IsIP, Addr, Intf automatically
// filled based on 's' argument parsing resut. If 's' is parsed to netip.Addr, then
// attribute 'Addr' set up to netip.Addr parsed, 'IsIP' set up to true. Otherwise,
// Intf set up to 's', IsIP set up to false
func newNextHop(s string) *nextHop {
	if v , err := netip.ParseAddr(s); err != nil {
		return &nextHop{IsIP: false, Intf: s}
	} else {
		return &nextHop{IsIP: true, Addr: v}
	}
}

// Needed for NXOS next-hops, where via part is always IP, regardless of route
// type (directly connected, local, etc.)
func (nh *nextHop) setIntf (s string) {
	nh.Intf = s
	nh.IsIP = false
	nh.Addr = netip.Addr{}
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
type RoutingTable struct{
	Table string
	Routes []*route
	NH map[uint64]*nextHop
}

// Constructor for Routing table. Default table name is 'default'
func newRoutingTable() *RoutingTable {
	return &RoutingTable{
		Table: "default",
		Routes: make([]*route, 0),
		NH: make(map[uint64]*nextHop),
	}
}

func (rt *RoutingTable) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("Table: %s\n", rt.Table))
	b.WriteString(fmt.Sprintf("Routes: %d\n", rt.RoutesCount()))
	for _,v := range rt.Routes {
		b.WriteString(v.String() + "\n")
	}
	b.WriteString(fmt.Sprintf("Next-Hops: %d\n", rt.NHCount()))
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

// FindRoutes func return number of routes found, channel of *route objects, which contain "ip"
// specified and error if any. Routes put in channel are ordered based on prefix lenght, 
// starting from more specific. If "all" flag is specified, func return all matched routes,
// otherwise only best match returned
func (rt *RoutingTable) FindRoutes(ip string, all bool) (int, <-chan *route, error) {
	count := 0
	out := make(chan *route)
	parsedIp, err := netip.ParseAddr(ip)
	if err != nil {
		close(out)
		return count, out, err
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
		return count, out, nil
	}
	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i].Network.Bits() > indexes[j].Network.Bits()
	})

	// for case when we need to return only best route
	if !all {
		indexes = indexes[:1]
	}
	count = len(indexes)
	go func(){
		defer close(out)
		for _,v := range indexes {
			out <- v
		}
	}()
	return count, out, nil
}

// FindRoutesByNH func finds all routes with specified nexthop.
// Returns number of routes found, and channel with those *route objects.
func (rt *RoutingTable) FindRoutesByNH(n string) (int, <-chan *route) {
	count := 0
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
		return count, out
	}
	count = len(res)
	go func(){
		defer close(out)
		for _,v := range res {
			out <- v
		}
	}()
	return count, out
}

// FindUniqNexthop func finds all unique *next-hop objects. 
// Returns number if next-hops found, and channel with those *next-hop objects.
// "ipOnly" flag gives ability to specify whether we need to get only NextHops with IP address
func (rt *RoutingTable) FindUniqNexthops (ipOnly bool) (int, <-chan *nextHop) {
	count := 0
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
	count = len(nhList)
	go func(){
		defer close(out)
		for _,v := range nhList {
			out <-v
		}
	}()
	return count, out
}