package main

import (
	"log"
	"net/netip"
	"os"
	"regexp"
	"strings"
	"fmt"
)


const (
	CONN_ROUTE_REGEXP = `^(\w\*? ?\w?\w?) +(\S+) is directly connected, (\S+)`
	LINE_BREAK_REGEXP = `\[.*] via (\S+)$`
	REGULAG_ROUTE_REGEXP = `^(\w\*? ?\w?\w?) +(\S+) [/\d\[\]]+ via (\S+)`
)

var text = []string{
	"O*E1  0.0.0.0/0 [110/101] via 192.168.199.18, 3w5d, Vlan14",
	"	   10.0.0.0/8 is variably subnetted, 7 subnets, 3 masks",
	"O E2     10.10.10.0/24 [110/20] via 192.168.199.18, 1w3d, Vlan14",
	"L        172.17.60.1/32 is directly connected, Vlan360",
	"O        172.17.61.0/24 [110/41] via 192.168.199.35, 1w5d, Vlan889",
	"					     [110/41] via 192.168.199.34, 1w5d, Vlan889",
}

var (
	InfoLogger  *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

var connRouteComp = regexp.MustCompile(CONN_ROUTE_REGEXP)
var lineBreakComp = regexp.MustCompile(LINE_BREAK_REGEXP)

var AllRoutes = Routes{}


func main() {

	var route = NewRoute()

	for _, line := range text {
		if strings.Contains(line, "is variably subnetted") {
			continue

		//case for directly connected route
		} else if strings.Contains(line, "is directly connected") {
			matches := connRouteComp.FindStringSubmatch(line)
			rtype := strings.TrimSpace(matches[1])
			prefix, err := netip.ParsePrefix(matches[2])
			if err != nil {
				WarnLogger.Printf("Cannot parse prefix from string %s", matches[2])
			}
			
			route.Type = rtype
			route.Network = prefix
			nh := NewNextHop(matches[3])
			route.AddNextHop(nh)
			AllRoutes.Add(route)
		
		} else if m,_ :=regexp.MatchString(REGULAG_ROUTE_REGEXP, line); m {
			matches := regexp.MustCompile(REGULAG_ROUTE_REGEXP).FindStringSubmatch(line)
			rtype := strings.TrimSpace(matches[1])
			prefix, err := netip.ParsePrefix(matches[2])
			if err != nil {
				WarnLogger.Printf("Cannot parse prefix from string %s", matches[2])
			}
			route.Type = rtype
			route.Network = prefix
			nh := NewNextHop(matches[3])
			route.AddNextHop(nh)
			AllRoutes.Add(route)
			
		} else if strings.HasPrefix(line, "[") {
			matches := lineBreakComp.FindStringSubmatch(line)
			nh := NewNextHop(matches[1])
			AllRoutes.GetLast().AddNextHop(nh)
		}
	}
	fmt.Println(AllRoutes)

}