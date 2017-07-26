package go_servers

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

const (
	// The tests will run a test server on this port.
	port    = 9654
	timeOut = 500 * time.Millisecond
)

func createListener(sleep time.Duration) (*http.Server, net.Listener, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		time.Sleep(sleep)
		rw.WriteHeader(http.StatusOK)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	return server, l, err
}

func TestListenerConnectionLimitInvalidParams(t *testing.T) {
	s, _, err := createListener(timeOut)
	if err != nil {
		t.Fatal(err)
	}
	server := &ConnectionLimitServer{
		server:      s,
		ListenLimit: -1,
	}

	err = server.ListenAndServe()

	if err == nil {
		t.Error(err)
	}

	server.server.Close()
}

func TestListenerConnectionLimitValidParams(t *testing.T) {
	srv, _, err := createListener(timeOut)
	if err != nil {
		t.Fatal(err)
	}

	s := &ConnectionLimitServer{
		server:      srv,
		ListenLimit: -1,
	}

	err = s.ListenAndServe()

	if err != nil {
		t.Error(err)
	}

	s.server.Close()
}
