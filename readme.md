# Cisco ip-route parser

Sometimes it's problematic to find actual route in routing table to match certain IP address (because of many matching routes with different masks, summarisation, etc.), espetially when you have only output from "show ip route" command, without live iteraction with router.
This app is intended for working with cisco routing tables. It parses routng table, and do actual route lookup based on different conditions (classic longest-match lookup, lookup based on next-hop value, etc.).
The following functionality is available:
- route lookup based on entered IP address
- route lookup based on entered next-hop value (either IP or interface)
- next-hop analysis (list of all unique next-hops)
  
App parses "show ip route" output from text file, and provides interactive menu with choises to user. IOS, IOS-XE, IOS-XR, NXOS "show ip route" outputs are supported.
___
## Usage
Upon startup, app looks for file to open and OS-family string (need to be specified with flags). The flag semantic is as follow

| Flag | Allowed values | Description |
| ---- | -------------- | ------------|
| -i  | file name | "show ip route" file to open |
| -os  | ios, nxos | OS-family, output is taken from. Use "nxos" for NXOS platforms, for all other use "ios" |

___
## Examples
Let's suppose we have the following routing table:
```
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
```
Upon successfull opening, will have output 
```
INFO: 2023/11/02 15:00:55 main.go:20: Starting...
INFO: 2023/11/02 15:00:55 main.go:35: Parsing routes...
INFO: 2023/11/02 15:00:55 main.go:38: Parsing routes done, found 14 routes, 10 unique nexthops

======================================
Possible values for selection:
1 - do route lookup based on entered IP
2 - do route lookup based on entered Next-hop
3 - print list of all unique Next-hops
8 - print raw routingTable object (for debug)
9 - exit the program
======================================
Enter your choise:
```
For next-hop analysis, we enter number 3. The output is below:
```
Found 10 unique nexthops:
{NextHop: 192.168.19.35}
{NextHop: 192.168.19.34}
{NextHop: Vlan8}
{NextHop: Null0}
{NextHop: Loopback102}
{NextHop: Port-channel2.21}
{NextHop: Port-channel2.32}
{NextHop: Port-channel1.189}
{NextHop: Loopback103}
{NextHop: 193.1.2.120}
```
Let's suppose we need to find all routes with next-hop "193.1.2.120". We enter 2
```
Enter your choise:
2
Enter Next-hop value, either IP or interface format accepted:
193.1.2.120
Found 2 routes:
S route to 193.1.3.2/32 network via [{NextHop: 193.1.2.120}]
S route to 193.1.2.0/24 network via [{NextHop: 193.1.2.120}]
```
Now let's find the route to 33.33.33.33 address. We enter 1
```
Enter your choise:
1
Enter IP:
33.33.33.33
Found 2 routes:
S route to 193.1.3.2/32 network via [{NextHop: 193.1.2.120}]
S route to 193.1.2.0/24 network via [{NextHop: 193.1.2.120}]
```
We got all matched routes, sorted from more specific to less specific.