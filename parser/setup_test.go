package parser

import (
	"testing"
	"os"
)

//var to store all parsed routes
var iosRoutes *RoutingTable
var nxosRoutes *RoutingTable

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(){
	infoLogger.Println("Setting up testing routing tables...")

	iosFile, err := os.Open("./testdata/iosRoute01.txt")
	if err != nil {
		errorLogger.Fatalf("Can not open IOSfile because of: %q", err)
	}
	defer iosFile.Close()
	nxosFile, err := os.Open("./testdata/nxosRoute01.txt")
	if err != nil {
		errorLogger.Fatalf("Can not open NXOSfile because of: %q", err)
	}
	defer iosFile.Close()
	defer nxosFile.Close()
	iosTableSource := NewTableSource("ios", iosFile)
	iosRoutes = iosTableSource.Parse()
	nxosTableSource := NewTableSource("nxos", nxosFile)
	nxosRoutes = nxosTableSource.Parse()
	
	infoLogger.Println("testing routing tables setup done")
}

func teardown() {
	infoLogger.Println("Testing finished")
}