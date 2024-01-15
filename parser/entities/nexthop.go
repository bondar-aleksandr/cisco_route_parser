package entities

import (
	"fmt"
	"log"
	"net/netip"
	"os"
	"github.com/mitchellh/hashstructure/v2"
)

var (
	InfoLogger  *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// Type to describe nexthop entity
type NextHop struct {
	IsIP bool
	Addr netip.Addr `hash:"string"`
	Intf string
}

// Returns pointer to new nextHop object. Attributes IsIP, Addr, Intf automatically
// filled based on 's' argument parsing resut. If 's' is parsed to netip.Addr, then
// attribute 'Addr' set up to netip.Addr parsed, 'IsIP' set up to true. Otherwise,
// Intf set up to 's', IsIP set up to false
func NewNextHop(s string) *NextHop {
	if v , err := netip.ParseAddr(s); err != nil {
		return &NextHop{IsIP: false, Intf: s}
	} else {
		return &NextHop{IsIP: true, Addr: v}
	}
}

// Needed for NXOS next-hops, where via part is always IP, regardless of route
// type (directly connected, local, etc.)
func (nh *NextHop) SetIntf (s string) {
	nh.Intf = s
	nh.IsIP = false
	nh.Addr = netip.Addr{}
}

func (nh *NextHop) String() string {
	if nh.IsIP {
		return fmt.Sprintf("{NextHop: %s}", nh.Addr)
	}
	return fmt.Sprintf("{NextHop: %s}", nh.Intf)
}

func (nh *NextHop) getHash() (uint64) {
	hash, err := hashstructure.Hash(nh, hashstructure.FormatV2, nil)
	if err != nil {
		ErrorLogger.Fatalf("Cannot compute hash from nexthop %s due to: %q", nh.String(), err)
	}
	return hash
}
