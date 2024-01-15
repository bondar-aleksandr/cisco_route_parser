package main

import (
	"flag"
	"os"
	"runtime/debug"
	"path/filepath"
	"github.com/bondar-aleksandr/cisco_route_parser/parser"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.SugaredLogger
)

func main() {

	// logger = zap.Must(zap.NewExample(), nil).Sugar()
	logger = createLogger().Sugar()
	defer logger.Sync()

	logger.Infof("Starting...")

	var iFileName = flag.String("i", "", "input 'ip route' filename to parse data from")
	var platform = flag.String("os", "", "OS family for the specified 'ip route' file. Allowed values are 'ios', 'nxos'")
	if len(os.Args) < 2 {
		logger.Fatalf("No input data provided, use -h flag for help. Exiting...")
	}
	flag.Parse()

	iFile, err := os.Open(*iFileName)
	if err != nil {
		logger.Fatalf("Can not open file %q because of: %q", *iFileName, err)
	}
	defer iFile.Close()

	logger.Infof("Parsing routes...")
	tableSource, err := parser.NewTableSource(*platform, iFile)
	if err != nil {
		logger.Errorf("got error: %w", err)
		os.Exit(1)
	}
	allRoutes := tableSource.Parse()
	logger.Infof("Parsing routes done, found %d routes, %d unique nexthops", allRoutes.RoutesCount(), allRoutes.NHCount())
	Menu(allRoutes)
}

func createLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	cwd, _ := os.Getwd()
	buildInfo, _ := debug.ReadBuildInfo()

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
        Development:       false,
        DisableCaller:     false,
        DisableStacktrace: false,
        Sampling:          nil,
        Encoding:          "json",
        EncoderConfig:     encoderCfg,
        OutputPaths: []string{
            "stderr",
			filepath.Join(cwd, "app.log"),
        },
        ErrorOutputPaths: []string{
            "stderr",
        },
        InitialFields: map[string]interface{}{
            "pid": os.Getpid(),
			"go version": buildInfo.GoVersion,
        },
	}
	return zap.Must(config.Build())
}