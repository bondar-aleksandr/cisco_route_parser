package main

import (
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

	// BuildRoutesCache(iFile, allRoutes)
	InfoLogger.Println("Parsing routes...")
	allRoutes = ParseRoute(iFile)
	InfoLogger.Printf("Parsing routes done, found %d routes, %d unique nexthops", allRoutes.Amount(), len(allNH))
	Menu()
}