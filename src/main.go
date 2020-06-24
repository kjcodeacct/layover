package main

import (
	"flag"
	"fmt"
	"layover/tcp"
	"layover/udp"
	"os"

	"go.uber.org/zap"
)

const helpPage = `
# Layover 
single port proxy process
<https://github.com/kjcodeacct/layover>

# Env Variables
* LAYOVER_PROXYHOST - default:"0.0.0.0"
	* the host layover is proxying from, unless specifying to a different host machine uses the default

* LAYOVER_PROXYPORT - required:true
	* the port layover is proxying *FROM*
    * this is *typically* the port not in the container system

* LAYOVER_PROTOCOL - default:"tcp"
	* IP protocol used by the specified port
    * options available 
        * "tcp"
        * "udp"

* LAYOVER_SERVEPORT default - default:"8080"
	* the port layover is proxying *TO* and is serving
    * if running in a container typically does *not* need to be specified

* LAYOVER_DEBUGMODE default - "0"
	* options available
		* 0 - off
		* 1 - basic logging of IP connecting and warnings
		* 2 - full logging including data (please don't use in production)

* LAYOVER_LOGDIR
	* directory to place logfiles created by enabling the LAYOVER_DEBUGMODE
`

var log *zap.SugaredLogger

var showHelp = flag.Bool("help", false, "show the help page")

func setup() {

	flag.Parse()

	zapLogger, _ := zap.NewProduction()

	defer zapLogger.Sync()
	log = zapLogger.Sugar()

	tcp.SetLog(log)
	udp.SetLog(log)

	if *showHelp {
		fmt.Println(helpPage)
		os.Exit(0)
	}

	err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	setup()

	if runtimeConfig.Protocol == udpMode {
		log.Info("Starting UDP server")
		udpListenAndServe()
	} else {
		log.Info("Starting TCP server")
		tcpListenAndServe()
	}
}
