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
	CONN_ROUTE_REGEXP = `^(\w\*? ?\w?\w?) +(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/?(\d{1,2})? is directly connected, ([^\,]+)`
	LINE_BREAK_REGEXP = `\[.*] via ([^\,]+)`
	REGULAG_ROUTE_REGEXP = `^(\w\*? ?\w?\w?) +(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/?(\d{1,2})? [/\d\[\]]+ via ([^\,]+)`
	COMMON_MASK_REGEXP = `/(\d{1,2})`
)

var text = []string{
	"B*    0.0.0.0/0 [20/0] via 212.26.135.74, 4w4d",
	"	   195.78.68.0/32 is subnetted, 1 subnets",
	"S        195.78.68.2 [1/0] via 195.78.69.120",
	"O*E1  0.0.0.0/0 [110/101] via 192.168.199.18, 3w5d, Vlan14",
	"	   10.0.0.0/8 is variably subnetted, 7 subnets, 3 masks",
	"O E2     10.10.10.0/24 [110/20] via 192.168.199.18, 1w3d, Vlan14",
	"C        195.78.69.112/28 is directly connected, Port-channel2.20",
	"L        195.78.69.119/32 is directly connected, Port-channel2.20",
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
var commonMaskComp = regexp.MustCompile(COMMON_MASK_REGEXP)

var AllRoutes = Routes{}


func main() {

	var commonMask string
	for _, line := range text {

		var route = NewRoute()
		
		//case where common mask specified for 
		if strings.Contains(line, "is subnetted") {
			matches := commonMaskComp.FindStringSubmatch(line)
			commonMask = matches[1]
			fmt.Println(commonMask)

		//case for directly connected route
		} else if strings.Contains(line, "is directly connected") {
			matches := connRouteComp.FindStringSubmatch(line)
			rtype := strings.TrimSpace(matches[1])
			pref := matches[2]
			mask := matches[3]
			if mask == "" {
				mask = commonMask
			}
			prefix, err := netip.ParsePrefix(fmt.Sprintf("%s/%s", pref, mask))
			if err != nil {
				WarnLogger.Printf("Cannot parse prefix from string %s/%s, skipping...", pref, mask)
				continue
			}
			route.Type = rtype
			route.Network = prefix
			nh := NewNextHop(matches[4])
			route.AddNextHop(nh)
			AllRoutes.Add(route)
		
		//case for regular route
		} else if m, _ :=regexp.MatchString(REGULAG_ROUTE_REGEXP, line); m {
			matches := regexp.MustCompile(REGULAG_ROUTE_REGEXP).FindStringSubmatch(line)
			rtype := strings.TrimSpace(matches[1])
			pref := matches[2]
			mask := matches[3]
			if mask == "" {
				mask = commonMask
			}
			prefix, err := netip.ParsePrefix(fmt.Sprintf("%s/%s", pref, mask))
			if err != nil {
				WarnLogger.Printf("Cannot parse prefix from string %s/%s, skipping...", pref, mask)
				continue
			}
			route.Type = rtype
			route.Network = prefix
			nh := NewNextHop(matches[4])
			route.AddNextHop(nh)
			AllRoutes.Add(route)
			
		//case for linebreak with via
		} else if strings.HasPrefix(strings.TrimSpace(line), "[") {
			matches := lineBreakComp.FindStringSubmatch(line)
			nh := NewNextHop(matches[1])
			AllRoutes.GetLast().AddNextHop(nh)
		}
	}
	fmt.Println(AllRoutes)

}