package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testMessage = "test-test-test"

func TestTCPListenAndServe(t *testing.T) {
	// adding an emtpy waitgroup just for the test
	var wg sync.WaitGroup

	go tcpListenAndServe(wg)

	// waiting for our proxy to startup
	time.Sleep(1 * time.Second)

	testConn := fmt.Sprintf("http://127.0.0.1:%d/test?msg=%s", ServePort, testMessage)
	t.Log("making request to", testConn)
	resp, err := http.Get(testConn)
	if err != nil {
		require.NoError(t, err)
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)

		receivedMsg := string(body)
		require.Equal(t, receivedMsg, testMessage)
	} else {
		t.Fatal("received status code", resp.StatusCode)
	}

}

func TestUDPListenAndServe(t *testing.T) {
	// adding an emtpy waitgroup just for the test
	var wg sync.WaitGroup

	go udpListenAndServe(wg)

	// waiting for our proxy to startup
	time.Sleep(1 * time.Second)

	testConn := fmt.Sprintf("127.0.0.1:%d", ServePort)
	t.Log("making request to", testConn)

	buff := make([]byte, len(testMessage))
	conn, err := net.Dial("udp", testConn)
	require.NoError(t, err)

	fmt.Fprintf(conn, testMessage)
	_, err = bufio.NewReader(conn).Read(buff)
	require.NoError(t, err)

	receivedMsg := string(buff)
	require.Equal(t, testMessage, receivedMsg)
}
