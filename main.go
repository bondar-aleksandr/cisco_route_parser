package main

import (
	"fmt"
	"log"
	"net/netip"
	"os"
	"strings"
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

func main() {
	r := strings.NewReader(strings.Join(text, "\n"))
	Routes := ParseRoute(r) 
	userIP, _ := netip.ParseAddr("195.78.69.119")
	rs := Routes.FindRoutes(userIP)
	fmt.Println(rs)

	// ipA, _ := netip.ParseAddr("1.2.3.4")
	// ipB, _ := netip.ParsePrefix("1.2.3.0/24")
	// ipC, _ := netip.ParsePrefix("1.2.3.8/29")
	// ipD, _ := netip.ParsePrefix("0.0.0.0/0")
	// fmt.Println(ipB.Contains(ipA))
	// fmt.Println(ipC.Contains(ipA))
	// fmt.Println(ipD.Contains(ipA))
}