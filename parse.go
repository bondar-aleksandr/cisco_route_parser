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
	REGULAG_ROUTE_REGEXP = `^(\w\*? ?\w?\w?) +(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/?(\d{1,2})? [/\d\[\]]+ via ([^\,]+)`
	COMMON_MASK_REGEXP = `/(\d{1,2})`
)

var connRouteComp = regexp.MustCompile(CONN_ROUTE_REGEXP)
var lineBreakComp = regexp.MustCompile(LINE_BREAK_REGEXP)
var commonMaskComp = regexp.MustCompile(COMMON_MASK_REGEXP)


func ParseRoute(r io.Reader) *Routes {

	var AllRoutes = &Routes{}
	var commonMask string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		var route = NewRoute()
		
		//case where common mask specified for 
		if strings.Contains(line, "is subnetted") {
			matches := commonMaskComp.FindStringSubmatch(line)
			commonMask = matches[1]

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
		//for debug purposes
		} else {
			WarnLogger.Printf("Line is not matched against any rule. Line: %s\n", line)
		}
	}
	return AllRoutes
}