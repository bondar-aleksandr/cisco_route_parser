package parser

import (
	// "fmt"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
)

var ipRoute = `
Routing Table: INET-ACCESS
Codes: L - local, C - connected, S - static, R - RIP, M - mobile, B - BGP
	D - EIGRP, EX - EIGRP external, O - OSPF, IA - OSPF inter area 
	N1 - OSPF NSSA external type 1, N2 - OSPF NSSA external type 2
	E1 - OSPF external type 1, E2 - OSPF external type 2
	i - IS-IS, su - IS-IS summary, L1 - IS-IS level-1, L2 - IS-IS level-2
	ia - IS-IS inter area, * - candidate default, U - per-user static route
	o - ODR, P - periodic downloaded static route, H - NHRP, l - LISP
	a - application route
	+ - replicated route, % - next hop override, p - overrides from PfR

Gateway of last resort is 212.26.135.74 to network 0.0.0.0

	193.1.3.0/32 is subnetted, 1 subnets
S        193.1.3.2 [1/0] via 193.1.2.120
	193.1.2.0/24 is variably subnetted, 5 subnets, 4 masks
S        193.1.2.0/24 [1/0] via 193.1.2.120
C        193.1.2.112/28 is directly connected, Port-channel2.21
L        193.1.2.119/32 is directly connected, Port-channel2.21
C        193.1.2.128/26 is directly connected, Port-channel2.32
L        193.1.2.190/32 is directly connected, Port-channel2.32
	189.110.135.0/24 is variably subnetted, 2 subnets, 2 masks
C        189.110.135.72/29 is directly connected, Port-channel1.189
L        189.110.135.77/32 is directly connected, Port-channel1.189
O        172.31.10.0/24 [110/41] via 192.168.19.35, 1w5d, Vlan8
                        [110/41] via 192.168.19.34, 1w5d, Vlan8
C        192.168.19.32/29 is directly connected, Vlan8
L        192.168.19.33/32 is directly connected, Vlan8
    33.0.0.0/8 is variably subnetted, 3 subnets, 2 masks
O        33.33.33.0/24 is a summary, 00:00:14, Null0
C        33.33.33.33/32 is directly connected, Loopback102
C        33.33.33.44/32 is directly connected, Loopback103
`

//var to store all parsed routes
var allRoutes *RoutingTable

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(){
	infoLogger.Println("Initiating testing routing table...")
	r := strings.NewReader(ipRoute)
	tableSource := NewTableSource("ios", r)
	allRoutes = tableSource.Parse()
}

func teardown() {
	infoLogger.Println("Testing finished")
}

func Test_Parsing(t *testing.T) {

	t.Run("parsing, amount of routes and NHs", func(t *testing.T) {
		assert.Equal(t, 14, allRoutes.RoutesCount())
		assert.Equal(t, 10, allRoutes.NHCount())
	})
}

func Test_VRF(t *testing.T) {
	t.Run("check table name", func(t *testing.T){
		assert.Equal(t, "INET-ACCESS", allRoutes.Table)
	})
}

func Test_RouteLookup(t *testing.T) {
	t.Run("correct ip present in routing table, exact match", func(t *testing.T) {
		routes, err := allRoutes.FindRoutes("189.110.135.77", false)
		if err != nil {
			t.Errorf("Parsing ip failed")
		}
		res := []*route{}
		for r := range routes {
			res = append(res, r)
		}
		assert.Equal(t, []*route{allRoutes.getByNetwork("189.110.135.77/32")}, res)
	})
	t.Run("correct ip present in routing table, all matches, correct output order", func(t *testing.T) {
		routes, err := allRoutes.FindRoutes("189.110.135.77", true)
		if err != nil {
			t.Errorf("Parsing ip failed")
		}
		res := []*route{}
		for r := range routes {
			res = append(res, r)
		}
		assert.Equal(t, []*route{allRoutes.getByNetwork("189.110.135.77/32"), allRoutes.getByNetwork("189.110.135.72/29")}, res)
	})
	t.Run("correct ip subnet address in routing table, exact match", func(t *testing.T) {
		routes, err := allRoutes.FindRoutes("33.33.33.0", false)
		if err != nil {
			t.Errorf("Parsing ip failed")
		}
		res := []*route{}
		for r := range routes {
			res = append(res, r)
		}
		assert.Equal(t, []*route{allRoutes.getByNetwork("33.33.33.0/24")}, res)
	})
	t.Run("correct ip address in routing table, multiple NH", func(t *testing.T) {
		routes, err := allRoutes.FindRoutes("172.31.10.0", false)
		if err != nil {
			t.Errorf("Parsing ip failed")
		}
		res := []*route{}
		for r := range routes {
			res = append(res, r)
		}
		assert.Equal(t, []*route{allRoutes.getByNetwork("172.31.10.0/24")}, res)
		assert.Equal(t, 2, (res[0].nhCount()))
	})
	t.Run("correct ip not present in routing table", func(t *testing.T) {
		routes, err := allRoutes.FindRoutes("1.2.3.4", false)
		if err != nil {
			t.Errorf("Parsing ip failed")
		}
		res := []*route{}
		for r := range routes {
			res = append(res, r)
		}
		assert.Equal(t, []*route{}, res)
	})
	t.Run("incorrect ip", func(t *testing.T) {
		_, err := allRoutes.FindRoutes("1.2.323.4", false)
		assert.ErrorContains(t, err, "has value >255")		
	})
	t.Run("incorrect ip with symbols", func(t *testing.T) {
		_, err := allRoutes.FindRoutes("ab.c.d s", false)
		assert.ErrorContains(t, err, "unexpected character")		
	})
	t.Run("blank ip", func(t *testing.T) {
		_, err := allRoutes.FindRoutes("", false)
		assert.ErrorContains(t, err, "unable to parse IP")		
	})
}

func Test_FindByNexthop(t *testing.T) {
	t.Run("correct next hop(IP), route present in routing table", func(t *testing.T) {
		res := []*route{}
		routes := allRoutes.FindRoutesByNH("192.168.19.35")
		for r := range routes {
			res = append(res, r)
		}
		assert.Equal(t, []*route{allRoutes.getByNetwork("172.31.10.0/24")}, res)
	})
	t.Run("correct next hop(interface), route present in routing table", func(t *testing.T) {
		res := []*route{}
		routes := allRoutes.FindRoutesByNH("Port-channel2.21")
		for r := range routes {
			res = append(res, r)
		}
		assert.Equal(t, []*route{allRoutes.getByNetwork("193.1.2.112/28"), allRoutes.getByNetwork("193.1.2.119/32")}, res)
	})
	t.Run("incorrect next hop", func(t *testing.T) {
		res := []*route{}
		routes := allRoutes.FindRoutesByNH("dfkffdkjs ds ")
		for r := range routes {
			res = append(res, r)
		}
		assert.Equal(t, []*route{}, res)
	})
	t.Run("blank next hop", func(t *testing.T) {
		res := []*route{}
		routes := allRoutes.FindRoutesByNH("")
		for r := range routes {
			res = append(res, r)
		}
		assert.Equal(t, []*route{}, res)
	})
}