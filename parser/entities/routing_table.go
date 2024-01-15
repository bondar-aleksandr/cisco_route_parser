package entities

import (
	"net/netip"
	"strings"
	"fmt"
	"sort"
)


// Routing table type. Consists of *Route slice and *nextHop map. Only unique next-hops are stored.
// Next-hops stored in map, where keys are their hashes, values are hext-hops themselves.
type RoutingTable struct{
	Table string
	Routes []*Route
	NH map[uint64]*NextHop
}

// Constructor for Routing table. Default table name is 'default'
func NewRoutingTable() *RoutingTable {
	return &RoutingTable{
		Table: "default",
		Routes: make([]*Route, 0),
		NH: make(map[uint64]*NextHop),
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

func (rt *RoutingTable) AddRoute(r *Route) {
	r.ParentRT = rt
	rt.Routes = append(rt.Routes, r)
}

// Only unique values added
func (rt *RoutingTable) addNextHop(nh *NextHop) {
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

func (rt *RoutingTable) GetLast() *Route {
	return rt.Routes[rt.RoutesCount() - 1]
}

// For test purposes. It's assumed that there is only one route to the destination in routing table
func (rt *RoutingTable) getByNetwork(s string) *Route {
	netw, err := netip.ParsePrefix(s)
	if err != nil {
		ErrorLogger.Printf("Cannot parse ip %s", s)
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
func (rt *RoutingTable) FindRoutes(ip string, all bool) (int, <-chan *Route, error) {
	count := 0
	out := make(chan *Route)
	parsedIp, err := netip.ParseAddr(ip)
	if err != nil {
		close(out)
		return count, out, err
	}

	indexes := []*Route{}
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
func (rt *RoutingTable) FindRoutesByNH(n string) (int, <-chan *Route) {
	count := 0
	out := make(chan *Route)
	nh := NewNextHop(n)
	res := []*Route{}
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
func (rt *RoutingTable) FindUniqNexthops (ipOnly bool) (int, <-chan *NextHop) {
	count := 0
	out := make(chan *NextHop)
	nhList := []*NextHop{}
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