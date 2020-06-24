package main

import (
	"fmt"
	"layover/tcp"
	"layover/udp"
	"net"
	"time"
)

// tcpListenAndServe - loops and listens to all incoming connections on specified TCP port
func tcpListenAndServe() {

	localConnection := fmt.Sprintf("%s:%d", "127.0.0.1", runtimeConfig.ServePort)
	portListen, err := net.Listen(runtimeConfig.Protocol, localConnection)
	if err != nil {
		log.Fatal("Failed to start up on port ", runtimeConfig.ServePort)
	}

	proxyConnection := fmt.Sprintf("%s:%d", runtimeConfig.ProxyHost, runtimeConfig.ProxyPort)
	proxy, err := tcp.NewProxy(proxyConnection, 10*time.Millisecond)
	if err != nil {
		log.Fatal("Failed to open connection to ", proxyConnection)
	}

	if proxy == nil {
		log.Fatal("Failed to proxy to ", proxyConnection)
	}

	log.Info("Accepting and proxying connections on ", runtimeConfig.ServePort)

	for {
		conn, err := portListen.Accept()
		if err != nil {
			log.Error(err)
		}

		proxy.ServeTCP(conn.(*net.TCPConn))
	}

	// should never hit this, but leaving for possible debugging
	portListen.Close()
}

// udpListenAndServe - loops and listens to all incoming connections on specified UDP port
func udpListenAndServe() {

	localConnection := fmt.Sprintf("%s:%d", "127.0.0.1", runtimeConfig.ServePort)
	addr, err := net.ResolveUDPAddr("udp", localConnection)
	if err != nil {
		log.Fatal("failed to start up on port ", runtimeConfig.ServePort)
	}

	portListen, err := udp.Listen(runtimeConfig.Protocol, addr)
	if err != nil {
		log.Fatal("failed to start up on port ", runtimeConfig.ServePort)
	}

	proxyConnection := fmt.Sprintf("%s:%d", runtimeConfig.ProxyHost, runtimeConfig.ProxyPort)
	proxy, err := udp.NewProxy(proxyConnection)
	if err != nil {
		log.Fatal("Failed to proxy connection to ", proxyConnection)
	}

	log.Info("Accepting and proxying connections on ", runtimeConfig.ServePort)

	for {

		conn, err := portListen.Accept()
		if err != nil {
			log.Error(err)
		}

		proxy.ServeUDP(conn)
	}

	// should never hit this, but leaving for possible debugging
	portListen.Close()
}
