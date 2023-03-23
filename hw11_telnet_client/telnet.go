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
		return err
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
	_, err := io.Copy(c.conn, c.in)
	return err
}

func (c *Client) Receive() error {
	if c.conn == nil {
		return errors.New("no connect")
	}
	_, err := io.Copy(c.out, c.conn)
	return err
}
