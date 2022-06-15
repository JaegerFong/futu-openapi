package tcp

import (
	"errors"
	"io"
	"net"
	"sync"
)

type Handler interface {
	Handle() error
}

type Encoder interface {
	WriteTo(c net.Conn) error
}

type Decoder interface {
	ReadFrom(c net.Conn) (Handler, error)
}

type Conn struct {
	c  net.Conn
	wg sync.WaitGroup
	de Decoder
}

func newConn(c net.Conn, de Decoder) *Conn {
	cn := &Conn{
		c:  c,
		de: de,
	}

	go cn.recv()
	return cn
}

func (c *Conn) Send(e Encoder) error {
	return e.WriteTo(c.c)
}

func (c *Conn) Close() error {
	if err := c.c.Close(); err != nil {
		return err
	}
	c.wg.Wait()
	return nil
}

// 持续接受数据，开启goroutine处理
func (c *Conn) recv() {
	for {
		h, err := c.de.ReadFrom(c.c)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				// TODO: log
				return
			}

			continue
		}

		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			err := h.Handle()
			if err != nil {
				// TODO: log
			}
			// TODO: log
		}()

	}
}

func Dial(addr string, de Decoder) (*Conn, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return newConn(c, de), nil
}
