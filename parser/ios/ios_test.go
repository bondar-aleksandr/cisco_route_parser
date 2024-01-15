package ios

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"errors"
)

func Test_Parsing(t *testing.T) {

	t.Run("parsing, amount of routes and NHs", func(t *testing.T) {
		assert.Equal(t, 14, iosRoutes.RoutesCount())
		assert.Equal(t, 10, iosRoutes.NHCount())
	})
}

func Test_VRF(t *testing.T) {
	t.Run("check table name", func(t *testing.T){
		assert.Equal(t, "INET-ACCESS", iosRoutes.Table)
	})
}

func Test_RouteLookupIOS(t *testing.T) {
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
			ip: "189.110.135.77",
			allRoutes: false,
			count: 1,
			routes: []*route{iosRoutes.getByNetwork("189.110.135.77/32")},
			err: nil,
		},
		{
			name: "correct ip present in routing table, all matches, correct output order",
			ip: "189.110.135.77",
			allRoutes: true,
			count: 2,
			routes: []*route{iosRoutes.getByNetwork("189.110.135.77/32"), iosRoutes.getByNetwork("189.110.135.72/29")},
			err : nil,
		},
		{
			name: "correct ip address in routing table, multiple NH",
			tType: "nhCount",
			ip: "172.31.10.0",
			allRoutes: false,
			count: 1,
			nhCount: 2,
			routes: []*route{iosRoutes.getByNetwork("172.31.10.0/24")},
			err : nil,
		},
		{
			name: "correct ip subnet address in routing table, exact match",
			ip: "33.33.33.0",
			allRoutes: false,
			count: 1,
			routes: []*route{iosRoutes.getByNetwork("33.33.33.0/24")},
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
			n, routes, err := iosRoutes.FindRoutes(tt.ip, tt.allRoutes)
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

func Test_FindByNexthopIOS(t *testing.T) {
	data := []struct{
		name string
		nh string
		count int
		routes []*route
	}{
		{
			name: "correct next hop(IP), route present in routing table",
			nh: "192.168.19.35",
			count: 1,
			routes: []*route{iosRoutes.getByNetwork("172.31.10.0/24")},
		},
		{
			name: "correct next hop(interface), route present in routing table",
			nh: "Port-channel2.21",
			count: 2,
			routes: []*route{iosRoutes.getByNetwork("193.1.2.112/28"), iosRoutes.getByNetwork("193.1.2.119/32")},
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
			n, routes := iosRoutes.FindRoutesByNH(tt.nh)
			res := []*route{}
			for r := range routes {
				res = append(res, r)
			}
			assert.ElementsMatch(t, tt.routes, res)
			assert.Equal(t, tt.count, n)
		})
	}
}