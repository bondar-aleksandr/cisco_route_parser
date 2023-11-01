package parser

import (
	"fmt"
	"strings"
	"net/netip"
	"regexp"
	"bufio"
)

const (
	conn_route_regexp = `^(\w\*? ?\w?\w?) +(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/?(\d{1,2})? is directly connected,( \w+,)? (\S+)$`
	line_break_regexp = `\[.*] via ([^\,]+)`
	regular_route_regexp = `^(\w\*? ?\w?\w?) +(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/?(\d{1,2})? [/\d\[\]]+ via ([^\,]+)`
	summary_route_regexp =`(\w\*? ?\w?\w?) +(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/?(\d{1,2})? is a summary(, .+?)?([^\,, ]+)$`
	common_mask_regexp = `/(\d{1,2})`
	table_name_regexp = `Routing Table: (\S+)`
)

var connRouteComp = regexp.MustCompile(conn_route_regexp)
var lineBreakComp = regexp.MustCompile(line_break_regexp)
var commonMaskComp = regexp.MustCompile(common_mask_regexp)
var summaryRouteComp = regexp.MustCompile(summary_route_regexp)
var tableNameComp = regexp.MustCompile(table_name_regexp)


func parseRouteIOS(t *tableSource) *RoutingTable {

	var RT = newRoutingTable()
	var commonMask string

	scanner := bufio.NewScanner(t.Source)
	for scanner.Scan() {
		line := scanner.Text()
	
		// case where common mask specified
		// 1.0.0.0/24 is subnetted, 10 subnets
		if strings.Contains(line, "is subnetted") {
			matches := commonMaskComp.FindStringSubmatch(line)
			commonMask = matches[1]

		// case where vrf name specified
		// Routing Table: INET-ACCESS
		} else if strings.Contains(line, "Routing Table:") {
			matches := tableNameComp.FindStringSubmatch(line)
			table := matches[1]
			RT.Table = table

		// case for directly connected route, for example:
		// C        33.33.33.33/32 is directly connected, Loopback102
		} else if strings.Contains(line, "is directly connected") {
			matches := connRouteComp.FindStringSubmatch(line)
			route := routeCreate(matches, []int{1,2,3,5}, commonMask, RT)
			if route == nil {
				continue
			}
			RT.addRoute(route)

		// case for summary discard route, for example:
		// O        33.33.33.0/24 is a summary, 00:00:14, Null0
		} else if strings.Contains(line, "is a summary") {
			matches := summaryRouteComp.FindStringSubmatch(line)
			route := routeCreate(matches, []int{1,2,3,5}, commonMask, RT)
			if route == nil {
				continue
			}
			RT.addRoute(route)
		
		// case for regular route, for example:
		// O        172.17.10.0/24 [110/41] via 192.168.199.35, 1w5d, Vlan889
		} else if m, _ := regexp.MatchString(regular_route_regexp, line); m {
			matches := regexp.MustCompile(regular_route_regexp).FindStringSubmatch(line)
			route := routeCreate(matches, []int{1,2,3,4}, commonMask, RT)
			if route == nil {
				continue
			}
			RT.addRoute(route)
			
		// case for linebreak with via
		// [110/41] via 192.168.199.34, 1w5d, Vlan889
		} else if strings.HasPrefix(strings.TrimSpace(line), "[") {
			matches := lineBreakComp.FindStringSubmatch(line)
			nh := newNextHop(matches[1])
			RT.getLast().addNextHop(nh)

		// just not for log.Warn to be triggered
		// 33.0.0.0/8 is variably subnetted, 3 subnets, 2 masks
		} else if strings.Contains(line, "is variably subnetted") {
			continue

		//for debug purposes
		// } else {
		// 	warnLogger.Printf("Line is not matched against any rule. Line: %s\n", line)
		}
	}
	return RT
}

// routeCreate func creates *Route object from slice of strings (matches) and corresponding
// indexes (regGroup) for those strings in the slice
func routeCreate(matches []string, capGroup []int, commonMask string, rt *RoutingTable) *route {
	var route = newRoute(rt)

	rtypeIndex := capGroup[0]
	prefIndex := capGroup[1]
	maskIndex := capGroup[2]
	nhIndex := capGroup[3]

	rtype := strings.TrimSpace(matches[rtypeIndex])
	pref := matches[prefIndex]
	mask := matches[maskIndex]
	if mask == "" {
		mask = commonMask
	}
	prefix, err := netip.ParsePrefix(fmt.Sprintf("%s/%s", pref, mask))
	if err != nil {
		warnLogger.Printf("Cannot parse prefix from string %s/%s, skipping...", pref, mask)
		return nil
	}
	route.Type = rtype
	route.Network = prefix
	nh := newNextHop(matches[nhIndex])
	route.addNextHop(nh)
	return route
}