package cmd

import (
	"fmt"
	"layover/tcp"
	"layover/udp"
	"net"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(proxyCmd)

	rootCmd.PersistentFlags().StringVarP(&ProxyHost, "proxyhost", "", "127.0.0.1",
		"host layover is proxy from")
	rootCmd.PersistentFlags().IntVarP(&ProxyPort, "proxyport", "", 8080,
		"host layover is proxy from")
	rootCmd.PersistentFlags().StringVarP(&ServeHost, "servehost", "", "127.0.0.1",
		"host layover is proxy to")
	rootCmd.PersistentFlags().IntVarP(&ServePort, "serveport", "", 8081,
		"host layover is proxy from")
	rootCmd.PersistentFlags().IntVarP(&ProxyDebugMode, "debugmode", "", 0,
		"level of proxy debugging (0 - off, 1 - basic logging , 2 - full logging")
}

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Proxy a specified port from a designated host.",
	Long: `Proxy a specified port from a designated host.

Variables:
* LAYOVER_PROXYHOST - default:"127.0.0.1"
	* the host layover is proxying *FROM*, unless specifying to a different host machine uses the default

* LAYOVER_PROXYPORT - required:true
	* the port layover is proxying *FROM*
	* this is *typically* the port not in the container system

* LAYOVER_SERVEHOST default - default:"127.0.0.1"
	* the host layover is proxying *TO*, unless specifying to a different host machine uses the default
	* if running in a container typically does *not* need to be specified

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

Example:
LAYOVER_PROXYHOST:LAYOVER_PROXYPORT -> LAYOVER_SERVEHOST:LAYOVER_SERVEPORT
`,
	Run: func(cmd *cobra.Command, args []string) {

		bindFlags(cmd, viper.GetViper())

		var wg sync.WaitGroup

		wg.Add(1)
		go tcpListenAndServe(wg)

		wg.Add(1)
		go udpListenAndServe(wg)

		wg.Wait()
	},
}

// tcpListenAndServe - loops and listens to all incoming connections on specified TCP port
func tcpListenAndServe(wg sync.WaitGroup) {

	tcp.SetLog(log)

	serveConnection := fmt.Sprintf("%s:%d", ServeHost, ServePort)
	portListen, err := net.Listen("tcp", serveConnection)
	if err != nil {
		log.Fatal("Failed to start up on port ", ServePort)
	}

	proxyConnection := fmt.Sprintf("%s:%d", ProxyHost, ProxyPort)
	// timeout := 10*time.Millisecond
	var timeout time.Duration
	tcpProxy, err := tcp.NewProxy(proxyConnection, timeout)
	if err != nil {
		log.Fatal("Failed to open connection to ", proxyConnection)
	}

	if tcpProxy == nil {
		log.Fatal("Failed to proxy to ", proxyConnection)
	}

	log.Info("Accepting and proxying TCP connections from ", proxyConnection, " to ", serveConnection)

	for {
		conn, err := portListen.Accept()
		if err != nil {
			log.Error(err)
		}

		tcpProxy.ServeTCP(conn.(*net.TCPConn))
	}

	// should never hit this, but leaving for possible debugging
	portListen.Close()
	wg.Done()
}

// udpListenAndServe - loops and listens to all incoming connections on specified UDP port
func udpListenAndServe(wg sync.WaitGroup) {

	udp.SetLog(log)

	serveConnection := fmt.Sprintf("%s:%d", ServeHost, ServePort)
	addr, err := net.ResolveUDPAddr("udp", serveConnection)
	if err != nil {
		log.Fatal("failed to start up on port ", ServePort)
	}

	portListen, err := udp.Listen("udp", addr)
	if err != nil {
		log.Info(err)
		log.Fatal("failed to start up on port ", ServePort)
	}

	proxyConnection := fmt.Sprintf("%s:%d", ProxyHost, ProxyPort)
	udpProxy, err := udp.NewProxy(proxyConnection)
	if err != nil {
		log.Fatal("Failed to proxy connection to ", proxyConnection)
	}

	log.Info("Accepting and proxying UDP connections from ", proxyConnection, " to ", serveConnection)

	for {

		conn, err := portListen.Accept()
		if err != nil {
			log.Error(err)
		}

		udpProxy.ServeUDP(conn)
	}

	// should never hit this, but leaving for possible debugging
	portListen.Close()
	wg.Done()
}
