package cmd

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/kjcodeacct/layover/tcp"
	"github.com/kjcodeacct/layover/udp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// ProxyHost - host layover is proxying from
	ProxyHost string
	// ProxyPort - port layover is proxying from
	ProxyPort int
	// ServeHost - host layover is proxying to
	ServeHost string
	// ServePort - port layover is proxying to
	ServePort int
)

func init() {
	rootCmd.AddCommand(proxyCmd)

	proxyCmd.PersistentFlags().StringVarP(&ProxyHost, "proxyhost", "", "127.0.0.1",
		"host layover is proxying from")
	proxyCmd.PersistentFlags().IntVarP(&ProxyPort, "proxyport", "", 8080,
		"port layover is proxying from")
	proxyCmd.PersistentFlags().StringVarP(&ServeHost, "servehost", "", "0.0.0.0",
		"host layover is proxying to")
	proxyCmd.PersistentFlags().IntVarP(&ServePort, "serveport", "", 8081,
		"port layover is proxying to")
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

	tcp.SetLog(logger)

	serveConnection := fmt.Sprintf("%s:%d", ServeHost, ServePort)
	portListen, err := net.Listen("tcp", serveConnection)
	if err != nil {
		logger.Fatal("Failed to start up on port ", ServePort)
	}

	proxyConnection := fmt.Sprintf("%s:%d", ProxyHost, ProxyPort)
	// timeout := 10*time.Millisecond
	var timeout time.Duration
	tcpProxy, err := tcp.NewProxy(proxyConnection, timeout)
	if err != nil {
		logger.Fatal("Failed to open connection to ", proxyConnection)
	}

	if tcpProxy == nil {
		logger.Fatal("Failed to proxy to ", proxyConnection)
	}

	logger.Info("Accepting and proxying TCP connections from ", proxyConnection, " to ", serveConnection)

	for {
		conn, err := portListen.Accept()
		if err != nil {
			logger.Error(err)
		}

		tcpProxy.ServeTCP(conn.(*net.TCPConn))
	}

	// should never hit this, but leaving for possible debugging
	portListen.Close()
	wg.Done()
}

// udpListenAndServe - loops and listens to all incoming connections on specified UDP port
func udpListenAndServe(wg sync.WaitGroup) {

	udp.SetLog(logger)

	serveConnection := fmt.Sprintf("%s:%d", ServeHost, ServePort)
	addr, err := net.ResolveUDPAddr("udp", serveConnection)
	if err != nil {
		logger.Fatal("failed to start up on port ", ServePort)
	}

	portListen, err := udp.Listen("udp", addr)
	if err != nil {
		logger.Info(err)
		logger.Fatal("failed to start up on port ", ServePort)
	}

	proxyConnection := fmt.Sprintf("%s:%d", ProxyHost, ProxyPort)
	udpProxy, err := udp.NewProxy(proxyConnection)
	if err != nil {
		logger.Fatal("Failed to proxy connection to ", proxyConnection)
	}

	logger.Info("Accepting and proxying UDP connections from ", proxyConnection, " to ", serveConnection)

	for {

		conn, err := portListen.Accept()
		if err != nil {
			logger.Error(err)
		}

		udpProxy.ServeUDP(conn)
	}

	// should never hit this, but leaving for possible debugging
	portListen.Close()
	wg.Done()
}
