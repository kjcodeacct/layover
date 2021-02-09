package tcp

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"go.uber.org/zap"
)

var log *zap.SugaredLogger

// this is a very common, but benign error, i.e the proxied connection has closed after
// recieving its data and not sending type appropriate headers
const tcpIgnoreError = "use of closed network connection"

// SetLog shares the log used by the main process
func SetLog(logInstance *zap.SugaredLogger) {
	log = logInstance
}

// below the methods for 'Proxy' and  'WriteCloser' are a modification and generalized adaption of the code found here
// https://github.com/containous/traefik/blob/master/pkg/tcp/proxy.go

// WriteCloser describes a net.Conn with a CloseWrite method.
type WriteCloser interface {
	net.Conn
	// CloseWrite on a network connection, indicates that the issuer of the call
	// has terminated sending on that connection.
	// It corresponds to sending a FIN packet.
	CloseWrite() error
}

// Proxy forwards a TCP request to a TCP service.
type Proxy struct {
	target           *net.TCPAddr
	terminationDelay time.Duration
}

// NewProxy creates a new Proxy.
func NewProxy(address string, terminationDelay time.Duration) (*Proxy, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Proxy{target: tcpAddr, terminationDelay: terminationDelay}, nil
}

// ServeTCP forwards the connection to a service.
func (p *Proxy) ServeTCP(conn WriteCloser) {

	// needed because of e.g. server.trackedConnection
	defer conn.Close()

	connBackend, err := net.DialTCP("tcp", nil, p.target)
	if err != nil {
		errMsg := fmt.Sprintf("Error while connection to backend: %v", err)
		log.Warn(errMsg)
		return
	}

	// maybe not needed, but just in case
	defer connBackend.Close()

	errChan := make(chan error)
	go p.connCopy(conn, connBackend, errChan)
	go p.connCopy(connBackend, conn, errChan)

	err = <-errChan
	if err != nil {
		errMsg := fmt.Sprintf("Error during connection: %v", err)
		log.Warn(errMsg)
	}

	<-errChan
}

func (p Proxy) connCopy(dst, src WriteCloser, errCh chan error) {
	_, err := io.Copy(dst, src)
	errCh <- err

	errClose := dst.CloseWrite()
	if errClose != nil {
		if !strings.Contains(errClose.Error(), tcpIgnoreError) {
			errMsg := fmt.Sprintf("Error while terminating connection: %v", errClose)
			log.Warn(errMsg)
		}
		return
	}

	if p.terminationDelay >= 0 {
		err := dst.SetReadDeadline(time.Now().Add(p.terminationDelay))
		if err != nil {
			if !strings.Contains(err.Error(), tcpIgnoreError) {
				errMsg := fmt.Sprintf("Error while setting deadline: %v", err)
				log.Warn(errMsg)
			}
		}
	}
}
