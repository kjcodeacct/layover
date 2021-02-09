package cmd

import (
	"fmt"
	"layover/tcp"
	"layover/udp"
	"net"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proxy called")

		var wg sync.WaitGroup

		wg.Add(1)
		go tcpListenAndServe(wg)

		wg.Add(1)
		go udpListenAndServe(wg)

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// proxyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// proxyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// tcpListenAndServe - loops and listens to all incoming connections on specified TCP port
func tcpListenAndServe(wg sync.WaitGroup) {

	localConnection := fmt.Sprintf("%s:%d", "127.0.0.1", ServePort)
	portListen, err := net.Listen("tcp", localConnection)
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

	log.Info("Accepting and proxying TCP connections on ", ServePort)

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

	localConnection := fmt.Sprintf("%s:%d", "127.0.0.1", ServePort)
	addr, err := net.ResolveUDPAddr("udp", localConnection)
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

	log.Info("Accepting and proxying UDP connections on ", ServePort)

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
