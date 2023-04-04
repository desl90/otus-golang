package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const (
	TelnetClientConnect   = "...Connected to "
	TelnetClientClose     = "...EOF"
	TelnetClientWasClosed = "...Connection was closed by peer"
)

var ErrTelnetClientNotConnected = errors.New("not connected")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		conn:    nil,
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type telnetClient struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (tc *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(os.Stderr, TelnetClientConnect+tc.address)
	if err != nil {
		return err
	}

	tc.conn = conn

	return nil
}

func (tc *telnetClient) Close() (err error) {
	if tc.conn == nil {
		return
	}

	if err = tc.conn.Close(); err != nil {
		return err
	}

	_, err = fmt.Fprintln(os.Stderr, TelnetClientClose)
	if err != nil {
		return err
	}

	return nil
}

func (tc *telnetClient) Send() error {
	if tc.conn == nil {
		_, err := fmt.Fprintln(os.Stderr, TelnetClientWasClosed)
		if err != nil {
			return err
		}

		return ErrTelnetClientNotConnected
	}

	_, err := io.Copy(tc.conn, tc.in)

	return err
}

func (tc *telnetClient) Receive() error {
	if tc.conn == nil {
		_, err := fmt.Fprintln(os.Stderr, TelnetClientWasClosed)
		if err != nil {
			return err
		}

		return ErrTelnetClientNotConnected
	}

	_, err := io.Copy(tc.out, tc.conn)

	return err
}
