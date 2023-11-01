package parser

import (
	"bufio"
	"net/netip"
	"regexp"
	"strings"
	"slices"
)

const (
	table_name_regexp_NXOS = `IP Route Table for VRF "(\S+)"`
	route_string_NXOS = `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d\d?`
	via_string_NXOS = `\*via (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(%\w+)?)?,? ?(\S+?)?, [/\d\[\]]+, [.:\w]+, ([\w-]+,? ?[\w-]+)`
)

var tableNameComp_NXOS = regexp.MustCompile(table_name_regexp_NXOS)
var routeStringComp_NXOS = regexp.MustCompile(route_string_NXOS)
var viaStringComp_NXOS = regexp.MustCompile(via_string_NXOS)

var direct_routes_NXOS = []string{"direct", "local", "hsrp", "glbp", "vrrp"}

func parseRouteNXOS(t *tableSource) *RoutingTable {

	var RT = newRoutingTable()
	var route = newRoute(RT)

	scanner := bufio.NewScanner(t.Source)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Table for VRF ") {
			matches := tableNameComp_NXOS.FindStringSubmatch(line)
			table := matches[1]
			RT.Table = table

		// case for line where prefix specified
		// 0.0.0.0/0, ubest/mbest: 1/0
		} else if strings.Contains(line, "ubest/mbest") {
			matches := routeStringComp_NXOS.FindStringSubmatch(line)
			pref, err := netip.ParsePrefix(matches[0])
			if err != nil {
				warnLogger.Printf("Cannot parse prefix from string %s, skipping...", pref)
				continue
			}
			route.Network = pref
			RT.addRoute(route)

		// case for line where next-hop specified
		// *via 192.168.199.33, Vlan889, [110/41], 1w3d, ospf-10, intra
		// *via 192.168.255.252, Lo0, [0/0], 2w5d, direct
		} else if strings.Contains(line, "*via ") {
			matches := viaStringComp_NXOS.FindStringSubmatch(line)
			nhStr := matches[1]
			nh := newNextHop(nhStr)
			nhIntf := matches[3]
			rtype := matches[4]

			// for cases where IP after '*via' needs to be replaced with interface (connected routes),
			// or where there is no IP after '*via' (summary, discard routes, routes via p2p interfaces, etc.)
			if slices.Contains(direct_routes_NXOS, rtype) || nhStr == "" {
				nh.setIntf(nhIntf)
			}
			RT.getLast().addNextHop(nh)
			RT.getLast().Type = rtype

			// create a new route for next iteration, since route is pointer
			route = newRoute(RT)
			
		//for debug purposes
		// } else {
		// 	warnLogger.Printf("Line is not matched against any rule. Line: %s\n", line)
		}

	}
	return RT
}