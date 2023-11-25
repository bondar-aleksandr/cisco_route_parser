package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/bondar-aleksandr/cisco_route_parser/parser"
	pb "github.com/bondar-aleksandr/cisco_route_parser/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	InfoLogger  *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)
//var to store all parsed routes
var allRoutes *parser.RoutingTable

const serverAddr = "localhost:50051"
const chunkSize = 1400

func main() {

	InfoLogger.Println("Starting...")

	var iFileName = flag.String("i", "", "input 'ip route' filename to parse data from")
	var platform = flag.String("os", "", "OS family for the specified 'ip route' file. Allowed values are 'ios', 'nxos'")
	if len(os.Args) < 2 {
		ErrorLogger.Fatalf("No input data provided, use -h flag for help. Exiting...")
	}
	flag.Parse()

	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial %s, : %v\n", serverAddr, err)
	}
	defer conn.Close()
	c := pb.NewRouteParserClient(conn)
	ctx := context.Background()
	Upload(ctx, c, iFileName, platform)

	

	// InfoLogger.Println("Parsing routes...")
	// tableSource := parser.NewTableSource(*platform, iFile)
	// allRoutes = tableSource.Parse()
	// InfoLogger.Printf("Parsing routes done, found %d routes, %d unique nexthops", allRoutes.RoutesCount(), allRoutes.NHCount())
	// Menu()
}