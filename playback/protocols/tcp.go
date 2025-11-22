package protocols

import (
	"io"
	"net"
	"strings"
)

type TCPProtocol struct{}

func (t *TCPProtocol) Connect(server string) (io.ReadCloser, error) {
	server = strings.TrimPrefix(server, "tcp://")
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(true)
		tcpConn.SetKeepAlive(true)
	}

	return conn, nil
}

func (t *TCPProtocol) Name() string {
	return "tcp"
}
