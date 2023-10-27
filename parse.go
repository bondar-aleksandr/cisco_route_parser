package main

import (
	"fmt"
	"strings"
	"net/netip"
	"regexp"
	"io"
	"bufio"
)

const (
	CONN_ROUTE_REGEXP = `^(\w\*? ?\w?\w?) +(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/?(\d{1,2})? is directly connected, ([^\,]+)`
	LINE_BREAK_REGEXP = `\[.*] via ([^\,]+)`
	REGULAR_ROUTE_REGEXP = `^(\w\*? ?\w?\w?) +(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/?(\d{1,2})? [/\d\[\]]+ via ([^\,]+)`
	SUMMARY_ROUTE_REGEXP =`(\w\*? ?\w?\w?) +(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/?(\d{1,2})? is a summary(, .+?)?([^\,, ]+)$`
	COMMON_MASK_REGEXP = `/(\d{1,2})`
)

var connRouteComp = regexp.MustCompile(CONN_ROUTE_REGEXP)
var lineBreakComp = regexp.MustCompile(LINE_BREAK_REGEXP)
var commonMaskComp = regexp.MustCompile(COMMON_MASK_REGEXP)
var summaryRouteComp = regexp.MustCompile(SUMMARY_ROUTE_REGEXP)


func ParseRoute(r io.Reader) *Routes {

	var AllRoutes = &Routes{}
	var commonMask string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		var route = NewRoute()
		
		// case where common mask specified
		// 1.0.0.0/24 is subnetted, 10 subnets
		if strings.Contains(line, "is subnetted") {
			matches := commonMaskComp.FindStringSubmatch(line)
			commonMask = matches[1]

		// case for directly connected route, for example:
		// C        33.33.33.33/32 is directly connected, Loopback102
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
			addNhToCache(nh)
			route.AddNextHop(nh.GetHash())
			AllRoutes.Add(route)

		// case for summary discard route, for example:
		// O        33.33.33.0/24 is a summary, 00:00:14, Null0
		} else if strings.Contains(line, "is a summary") {
			matches := summaryRouteComp.FindStringSubmatch(line)
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
			nh := NewNextHop(matches[5])
			addNhToCache(nh)
			AllRoutes.Add(route)
		
		// case for regular route, for example:
		// O        172.17.10.0/24 [110/41] via 192.168.199.35, 1w5d, Vlan889
		} else if m, _ :=regexp.MatchString(REGULAR_ROUTE_REGEXP, line); m {
			matches := regexp.MustCompile(REGULAR_ROUTE_REGEXP).FindStringSubmatch(line)
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
			addNhToCache(nh)
			route.AddNextHop(nh.GetHash())
			AllRoutes.Add(route)
			
		// case for linebreak with via
		// [110/41] via 192.168.199.34, 1w5d, Vlan889
		} else if strings.HasPrefix(strings.TrimSpace(line), "[") {
			matches := lineBreakComp.FindStringSubmatch(line)
			nh := NewNextHop(matches[1])
			addNhToCache(nh)
			AllRoutes.GetLast().AddNextHop(nh.GetHash())

		// just not for log.Warn to be triggered
		// 33.0.0.0/8 is variably subnetted, 3 subnets, 2 masks
		} else if strings.Contains(line, "is variably subnetted") {
			continue

		//for debug purposes
		// } else {
		// 	WarnLogger.Printf("Line is not matched against any rule. Line: %s\n", line)
		}
	}
	return AllRoutes
}

// buildNHCache builds NH cache as map, where keys are hashes and values are *nextHop
func addNhToCache(nh *nextHop) {
	if _, ok := allNH[nh.GetHash()]; ok {
		return
	}
	allNH[nh.GetHash()] = nh
}

// func routeCreate(matches []string) {

// }