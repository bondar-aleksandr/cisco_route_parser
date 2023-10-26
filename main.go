package main

import (
	"fmt"
	"log"
	"os"
	"flag"
)

var (
	InfoLogger  *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)
//var to store all parsed routes
var allRoutes *Routes

//var to store all parsed next-hops
var allNH = make(map[uint64]*nextHop)

func main() {

	InfoLogger.Println("Starting...")

	var iFileName = flag.String("i", "", "input 'ip route' file to parse data from")
	if len(os.Args) < 1 {
		ErrorLogger.Fatalf("No input data provided, use -h flag for help. Exiting...")
	}
	flag.Parse()

	iFile, err := os.Open(*iFileName)
	if err != nil {
		ErrorLogger.Fatalf("Can not open file %q because of: %q", *iFileName, err)
	}
	defer iFile.Close()

	// r := strings.NewReader(strings.Join(text, "\n"))
	buildRoutesCache(iFile)

	fmt.Println("Type IP value for route lookup:")
	var userIP string
	fmt.Scanln(&userIP)

	rs,_ := allRoutes.FindRoutes(userIP, false)
	for s := range rs {
		fmt.Println(s)
	}
	// for nh := range allRoutes.FindUniqNexthops(false) {
	// 	fmt.Println(nh)
	// }

	// var userNh string
	// fmt.Println("Type Next-hop value:")
	// fmt.Scanln(&userNh)

	// for rByNH := range allRoutes.FindRoutesByNH(userNh) {
	// 	fmt.Println(rByNH)
	// }
	// r := allRoutes.getByNetwork("192.168.99.0/24")
	// fmt.Println(r)


}