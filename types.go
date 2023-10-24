package main

import (
	"net/netip"
	"fmt"
	"sort"
)

//Single route entity
type Route struct {
	Network netip.Prefix
	Type string
	NHList []*nextHop
}

func (r *Route) String() string {
	return fmt.Sprintf("%s route to %s network via %v", r.Type, r.Network.String(), r.NHList)
}

//Constructor
func NewRoute() *Route {
	return &Route{
		NHList: make([]*nextHop, 0),
	}
}

func (r *Route) AddNextHop(nh *nextHop) {
	r.NHList = append(r.NHList, nh)
}


//Nexthop entity
type nextHop struct {
	isIP bool
	ip netip.Addr
	intf string
}

func NewNextHop(s string) *nextHop {
	if v , err := netip.ParseAddr(s); err != nil {
		return &nextHop{isIP: false, intf: s}
	} else {
		return &nextHop{isIP: true, ip: v}
	}
}

func (nh *nextHop) String() string {
	if nh.isIP {
		return fmt.Sprintf("NextHop: %s, isIP: %v", nh.ip, nh.isIP)
	}
	return fmt.Sprintf("NextHop: %s, isIP: %v", nh.intf, nh.isIP)
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

// FindRoutes func return slice of *Route objects, which contain "ip".
// Routes in slice are sorted based on prefix lenght, starting from more specific
func (r *Routes) FindRoutes(ip netip.Addr) []*Route {
	routes := []*Route{}
	for _, v := range r.Elements {
		if v.Network.Contains(ip) {
			routes = append(routes, v)
		}
	}
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Network.Bits() > routes[j].Network.Bits()
	})
	return routes
}