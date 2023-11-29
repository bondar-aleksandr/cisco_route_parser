package main

import (
	"context"
	"flag"
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)
//var to store all parsed routes
// var allRoutes *parser.RoutingTable
var c *ClientService

func main() {

	InfoLogger.Println("Starting...")

	var iFileName = flag.String("i", "", "input 'ip route' filename to parse data from")
	var platform = flag.String("os", "", "OS family for the specified 'ip route' file. Allowed values are 'ios', 'nxos'")
	if len(os.Args) < 2 {
		ErrorLogger.Fatalf("No input data provided, use -h flag for help. Exiting...")
	}
	flag.Parse()

	ctx := context.Background()
	c = NewClientService()
	
	defer func(){
		c.sessionClose(ctx)
		// c.Close()
	}()

	err := c.Upload(ctx, iFileName, platform)
	if err != nil {
		log.Fatalf("Failed to upload file to gRPC server due to: %v", err)
	}
	Menu(ctx)
}