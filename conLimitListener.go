package go_servers

import (
	"net"
	"sync"
)

// ConnLimitListener returns a Listener that accepts at most limit simultaneous
// connections from the provided Listener.
func ConnLimitListener(l net.Listener, limit int) net.Listener {
	return &connLimitListener{l, make(chan struct{}, limit)}
}

type connLimitListener struct {
	net.Listener
	connectionLimit chan struct{}
}

func (l *connLimitListener) acquireConnection() { l.connectionLimit <- struct{}{} }
func (l *connLimitListener) releaseConnection() { <-l.connectionLimit }

// Attempts to accept the connection as long as we haven't reached the
// connectionLimit
func (l *connLimitListener) Accept() (net.Conn, error) {
	l.acquireConnection()
	c, err := l.Listener.Accept()
	if err != nil {
		l.releaseConnection()
		return nil, err
	}

	return &limitListenerConn{Conn: c, releaseConnection: l.releaseConnection}, nil
}

type limitListenerConn struct {
	net.Conn
	once              sync.Once
	releaseConnection func()
}

// Closes the connection and releases the hold
func (l *limitListenerConn) Close() error {
	err := l.Conn.Close()
	l.once.Do(l.releaseConnection)
	return err
}
