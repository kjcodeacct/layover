package udp

import (
	"fmt"
	"io"
	"net"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

// SetLog shares the log used by the main process
func SetLog(logInstance *zap.SugaredLogger) {
	logger = logInstance
}

// Below the methods for 'udpProxy' and 'WriteCloser' are a modification and generalized adaption of the code found here
// https://github.com/containous/traefik/blob/master/pkg/udp/proxy.go

// Proxy is a reverse-proxy implementation of the udp handler interface.
type Proxy struct {
	target string
}

// NewProxy creates a new Proxy.
func NewProxy(address string) (*Proxy, error) {
	return &Proxy{target: address}, nil
}

// ServeUDP implements the Handler interface.
func (p *Proxy) ServeUDP(conn *Conn) {

	// needed because of e.g. server.trackedConnection
	defer conn.Close()

	connBackend, err := net.Dial("udp", p.target)
	if err != nil {
		errMsg := fmt.Sprintf("Error while connecting to backend: %v", err)
		logger.Warn(errMsg)
		return
	}

	// maybe not needed, but just in case
	defer connBackend.Close()

	errChan := make(chan error)
	go p.connCopy(conn, connBackend, errChan)
	go p.connCopy(connBackend, conn, errChan)

	err = <-errChan
	if err != nil {
		errMsg := fmt.Sprintf("Error while serving UDP: %v", err)
		logger.Warn(errMsg)
	}

	<-errChan
}

func (p Proxy) connCopy(dst io.WriteCloser, src io.Reader, errCh chan error) {
	_, err := io.Copy(dst, src)
	errCh <- err

	if err := dst.Close(); err != nil {
		errMsg := fmt.Sprintf("Error while terminating connection: %v", err)
		logger.Warn(errMsg)
	}
}
