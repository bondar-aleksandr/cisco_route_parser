package parser

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParsingNXOS(t *testing.T) {
	t.Run("parsing, amount of routes and NHs", func(t *testing.T) {
		assert.Equal(t, 135, nxosRoutes.RoutesCount())
		assert.Equal(t, 13, nxosRoutes.NHCount())
	})
}

func Test_vrfNXOS(t *testing.T) {
	t.Run("check table name", func(t *testing.T){
		assert.Equal(t, "default", nxosRoutes.Table)
	})
}

func Test_RouteLookupNXOS(t *testing.T) {
	data := []struct{
		name string
		tType string
		ip string
		allRoutes bool
		count int
		nhCount int
		routes []*route
		err any
	}{
		{
			name: "correct ip present in routing table, exact match",
			ip: "172.17.7.21",
			allRoutes: false,
			count: 1,
			routes: []*route{nxosRoutes.getByNetwork("172.17.7.0/24")},
			err: nil,
		},
		{
			name: "correct ip present in routing table, all matches, correct output order",
			ip: "10.224.100.1",
			allRoutes: true,
			count: 2,
			routes: []*route{nxosRoutes.getByNetwork("10.224.100.1/32"), nxosRoutes.getByNetwork("10.224.100.0/30")},
			err : nil,
		},
		{
			name: "correct ip address in routing table, multiple NH",
			tType: "nhCount",
			ip: "1.2.3.4",
			allRoutes: false,
			count: 1,
			nhCount: 2,
			routes: []*route{nxosRoutes.getByNetwork("1.2.3.4/32")},
			err : nil,
		},
		{
			name: "correct ip subnet address in routing table, exact match",
			ip: "172.17.2.0",
			allRoutes: false,
			count: 1,
			routes: []*route{nxosRoutes.getByNetwork("172.17.2.0/24")},
			err : nil,
		},
		{
			name: "correct ip not present in routing table",
			ip: "11.12.13.14",
			allRoutes: false,
			count: 0,
			routes: []*route{},
			err : nil,
		},
		{
			name: "incorrect IP",
			tType: "errCheck",
			ip: "11.12.313.14",
			allRoutes: false,
			count: 0,
			routes: []*route{},
			err : errors.New("xz"),
		},
		{
			name: "incorrect ip with symbols",
			tType: "errCheck",
			ip: "ab.c.d s2",
			allRoutes: false,
			count: 0,
			routes: []*route{},
			err : errors.New("xz"),
		},
		{
			name: "blank ip",
			tType: "errCheck",
			ip: "",
			allRoutes: false,
			count: 0,
			routes: []*route{},
			err : errors.New("xz"),
		},
	}
		
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			n, routes, err := nxosRoutes.FindRoutes(tt.ip, tt.allRoutes)
			res := []*route{}
			for r := range routes {
				res = append(res, r)
			}
			assert.Equal(t, tt.routes, res)
			assert.Equal(t, tt.count, n)
			if tt.tType == "nhCount" {
				assert.Equal(t, tt.nhCount, res[0].nhCount())
			} else if tt.tType == "errCheck" {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func Test_FindByNexthopNXOS(t *testing.T) {
	data := []struct{
		name string
		nh string
		count int
		routes []*route
	}{
		{
			name: "correct next hop(IP), route present in routing table",
			nh: "192.168.255.254",
			count: 1,
			routes: []*route{nxosRoutes.getByNetwork("1.2.3.4/32")},
		},
		{
			name: "correct next hop(interface), route present in routing table",
			nh: "Vlan889",
			count: 2,
			routes: []*route{nxosRoutes.getByNetwork("192.168.199.34/32"), nxosRoutes.getByNetwork("192.168.199.32/29")},
		},
		{
			name: "incorrect next hop",
			nh: "das dsda s123",
			count: 0,
			routes: []*route{},
		},
		{
			name: "blank next hop",
			nh: "",
			count: 0,
			routes: []*route{},
		},
	}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			n, routes := nxosRoutes.FindRoutesByNH(tt.nh)
			res := []*route{}
			for r := range routes {
				res = append(res, r)
			}
			assert.ElementsMatch(t, tt.routes, res)
			assert.Equal(t, tt.count, n)
		})
	}
}