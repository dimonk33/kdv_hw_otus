package main

import (
	"errors"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	client := Client{
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
	return &client
}

type Client struct {
	addr    string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *Client) Connect() error {
	var err error
	c.conn, err = net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	if err := c.in.Close(); err != nil {
		println(err)
	}

	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) Send() error {
	if c.conn == nil {
		return errors.New("no connect")
	}
	var count int
	var err error
	buf := make([]byte, 1024)
	for count == 0 {
		count, err = c.in.Read(buf)
		if err != nil {
			return err
		}
	}
	count, err = c.conn.Write(buf[:count])
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("no data sent")
	}
	return nil
}

func (c *Client) Receive() error {
	if c.conn == nil {
		return errors.New("no connect")
	}
	var count int
	var err error
	buf := make([]byte, 1024)
	for count == 0 {
		count, err = c.conn.Read(buf)
		if err != nil {
			return err
		}
	}
	count, err = c.out.Write(buf[:count])
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("no data write")
	}
	return nil
}
