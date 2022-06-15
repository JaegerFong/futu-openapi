package tcp

import (
	"bufio"
	"net"
	"testing"
)

func TestConnect(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Error(err)
	}

	data := make([]byte, 64)
	bufio.NewReader(conn).Read(data)
}
