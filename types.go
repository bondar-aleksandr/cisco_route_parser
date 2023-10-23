package main

import (
	"net/netip"
	"fmt"
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
		return fmt.Sprintf("NextHop: %s", nh.ip)
	}
	return fmt.Sprintf("NextHop: %s", nh.intf)
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