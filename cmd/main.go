package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/meta-node-blockchain/meta-node-mns/app"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
)

const (
	defaultConfigPath = "config.yaml"
	defaultLogLevel   = logger.FLAG_INFO
)

var (
	// flags
	CONFIG_FILE_PATH string
	LOG_LEVEL        int
)

func main() {

	flag.StringVar(&CONFIG_FILE_PATH, "config", defaultConfigPath, "Config path")
	flag.StringVar(&CONFIG_FILE_PATH, "c", defaultConfigPath, "Config path (shorthand)")

	flag.IntVar(&LOG_LEVEL, "log-level", defaultLogLevel, "Log level")
	flag.IntVar(&LOG_LEVEL, "ll", defaultLogLevel, "Log level (shorthand)")

	flag.Parse()

	app, err := app.NewApp(defaultConfigPath, LOG_LEVEL)
	if err != nil {
		panic(err)
	}
	go func() {
		app.Run()
	}()

	logger.Debug("Program run")
	sigs := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		close(done)
	}()
	<-done

}
