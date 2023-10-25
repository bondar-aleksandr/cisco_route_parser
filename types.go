package main

import (
	"fmt"
	"net/netip"
	"sort"
	"github.com/mitchellh/hashstructure/v2"
)

//Single route entity
type Route struct {
	Network netip.Prefix
	Type string
	NHList []uint64
}

func (r *Route) String() string {
	nhlist := []*nextHop{}
	for _,v := range r.NHList {
		nhlist = append(nhlist, allNH[v])
	}
	return fmt.Sprintf("%s route to %s network via %v", r.Type, r.Network.String(), nhlist)
}

//Constructor
func NewRoute() *Route {
	return &Route{
		NHList: make([]uint64, 0),
	}
}

func (r *Route) AddNextHop(nhHash uint64) {
	r.NHList = append(r.NHList, nhHash)
}

//Nexthop entity
type nextHop struct {
	IsIP bool
	Addr netip.Addr `hash:"string"`
	Intf string
}

func NewNextHop(s string) *nextHop {
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

func (nh *nextHop) GetHash() (uint64) {
	hash, err := hashstructure.Hash(nh, hashstructure.FormatV2, nil)
	if err != nil {
		ErrorLogger.Fatalf("Cannot compute hash from nexthop %s due to: %q", nh.String(), err)
	}
	return hash
}

//All routes
type Routes struct{
	Elements []*Route
}

func (r *Routes) Add(e *Route) {
	r.Elements = append(r.Elements, e)
}

func (r *Routes) Amount() int {
	return len(r.Elements)
}

func (r *Routes) GetLast() *Route {
	return r.Elements[r.Amount() - 1]
}

// FindRoutes func return channel of *Route objects, which contain "ip".
// Routes in slice are sorted based on prefix lenght, starting from more specific
func (r *Routes) FindRoutes(ip netip.Addr) <-chan *Route {
	out := make(chan *Route)
	indexes := []*Route{}
	for _, v := range r.Elements {
		if v.Network.Contains(ip) {
			indexes = append(indexes, v)
		}
	}
	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i].Network.Bits() > indexes[j].Network.Bits()
	})
	go func(){
		defer close(out)
		for _,v := range indexes {
			out <- v
		}
	}()
	return out
}

// FindRoutesByNH func finds all routes with specified nexthop.
// Result returned as a channel
func (r *Routes) FindRoutesByNH(n string) <-chan *Route {
	out := make(chan *Route)
	nh := NewNextHop(n)
	res := []*Route{}
	for _, route := range r.Elements {
		for _, v := range route.NHList {
			if v == nh.GetHash() {
				res = append(res, route)
			}
		}
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
func (r *Routes) FindUniqNexthops (ipOnly bool) <-chan *nextHop {
	out := make(chan *nextHop)
	nhList := []*nextHop{}
	for _, nh := range allNH {
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