package entities

import (
	"net/netip"
	"fmt"
)


// Type to describe route entity. 'ParentRT' attribute refers to *RoutingTable
// this route belongs. NHList attribute stores next-hop hashes, which are used by
// this route.
type Route struct {
	Network netip.Prefix
	Type string
	NHList []uint64
	ParentRT *RoutingTable
}

func (r *Route) String() string {
	nhlist := []*NextHop{}
	for _,v := range r.NHList {
		nhlist = append(nhlist, r.ParentRT.NH[v])
	}
	return fmt.Sprintf("%s route to %s network via %v", r.Type, r.Network.String(), nhlist)
}

// Returns pointer to new route object with 'ParentRT' attribute set up to rt value
func NewRoute(rt *RoutingTable) *Route {
	return &Route{
		NHList: make([]uint64, 0),
		ParentRT: rt,
	}
}

// to add nextHop to route. Whenever NH is added to the route object, 
// it's also added to RoutingTable object
func (r *Route) AddNextHop(nh *NextHop) {
	r.ParentRT.addNextHop(nh)
	r.NHList = append(r.NHList, nh.getHash())
}

// next-hops count
func (r *Route) NhCount() int {
	return len(r.NHList)
}