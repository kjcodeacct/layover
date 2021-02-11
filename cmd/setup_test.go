package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	var err error

	ProxyPort, err = testGetOpenPort()
	if err != nil {
		errMsg := fmt.Sprintf("failed to find open port %d: %s", ProxyPort, err.Error())
		log.Fatalln(errMsg)
	}

	log.Printf("testing proxy to port %d\n", ProxyPort)

	go testHTTPServer(ProxyPort)

	go testUDPServer(ProxyPort)

	ServePort, err = testGetOpenPort()
	if err != nil {
		errMsg := fmt.Sprintf("failed to find open port %d: %s", ServePort, err.Error())
		log.Fatalln(errMsg)
	}

	log.Printf("testing serve on port %d \n", ServePort)

	os.Exit(m.Run())
}

func testGetOpenPort() (int, error) {

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}

	openPort := listener.Addr().(*net.TCPAddr).Port

	listener.Close()

	return openPort, nil
}

func testHTTPServer(port int) {

	portStr := fmt.Sprintf(":%d", port)
	log.Printf("starting up on port %d\n", port)

	http.HandleFunc("/test", testHTTPHandler)

	http.ListenAndServe(portStr, nil)
}

func testHTTPHandler(w http.ResponseWriter, req *http.Request) {

	param := req.URL.Query().Get("msg")
	url := fmt.Sprintf("%v %v %v", req.Method, req.URL, req.Proto)
	log.Println(url)
	fmt.Fprintf(w, "%s", param)
}

func testUDPServer(port int) {

	udpAddr := &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("127.0.0.1"),
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalln("failed to start up on port ", ServePort)
	}

	// defaulting to 1024 bytes, MTU will probably stop anything over ~ 1400 bytes anyway
	// this solely for a test message, so it is a safe assumption
	buff := make([]byte, 1024)

	for {

		_, remoteAddr, err := conn.ReadFromUDP(buff)
		if err != nil {
			log.Fatalln(err)
		}

		testUDPHandler(conn, remoteAddr, buff)
	}
}

func testUDPHandler(conn *net.UDPConn, addr *net.UDPAddr, buff []byte) {
	_, err := conn.WriteToUDP(buff, addr)
	if err != nil {
		log.Fatalln(err)
	}
}
