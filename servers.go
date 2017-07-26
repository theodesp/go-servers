package go_servers

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

// Base server type.
type Server struct {
	*http.Server
}

// A server with a connection limit
type ConnectionLimitServer struct {
	*Server

	// Limit the number of outstanding requests
	ListenLimit int
}

// ListenAndServe is equivalent to http.Server.ListenAndServe
func (srv *ConnectionLimitServer) ListenAndServe() error {
	// Create the listener so we can control the connections coming.
	addr := srv.Server.Addr
	if addr == "" {
		addr = ":http"
	}
	conn, err := srv.newTCPListener(addr)
	if err != nil {
		return err
	}

	return srv.Serve(conn)
}

// Serve is equivalent to http.Server.Serve with a connection limit.
func (srv *ConnectionLimitServer) Serve(listener net.Listener) error {
	if srv.ListenLimit < 0 {
		return errors.New("Invalid Server Configuration: ListenLimit is negative")
	}
	if srv.ListenLimit != 0 {
		listener = ConnLimitListener(listener, srv.ListenLimit)
	}

	// Serve with graceful listener.
	// Execution blocks here until listener.Close() is called, above.
	err := srv.Server.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}

	// Graceful Shutdown
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Server.Shutdown(ctx)

	return err
}

func (srv *Server) newTCPListener(addr string) (net.Listener, error) {
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return conn, err
	}

	return conn, nil
}
