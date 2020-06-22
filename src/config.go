package main

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

const (
	DebugModeOff  = 0
	DebugModeMin  = 1
	DebugModeFull = 2

	tcpMode = "tcp"
	udpMode = "udp"

	processname = "layover"
)

type config struct {
	ProxyHost string
	ProxyPort int
	Protocol  string
	ServePort int
	DebugMode int
	LogDir    string
}

var runtimeConfig config

func readConfig() error {

	err := envconfig.Process(processname, &runtimeConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Running with config", runtimeConfig)

	if runtimeConfig.ProxyPort == 0 {
		log.Fatal("Must provide valid port to proxy")
	}

	if runtimeConfig.ProxyHost == "" {
		runtimeConfig.ProxyHost = "0.0.0.0"
	}

	if runtimeConfig.ServePort == 0 {
		runtimeConfig.ServePort = 8080
	}

	runtimeConfig.Protocol = strings.ToLower(runtimeConfig.Protocol)

	if runtimeConfig.Protocol == "" {
		runtimeConfig.Protocol = tcpMode
	}

	if runtimeConfig.Protocol != tcpMode {
		if runtimeConfig.Protocol != udpMode {
			errMsg := fmt.Sprintf("Connection mode of '%s' must be tcp or udp", runtimeConfig.Protocol)
			log.Fatal(errMsg)
		}
	}

	return nil
}
