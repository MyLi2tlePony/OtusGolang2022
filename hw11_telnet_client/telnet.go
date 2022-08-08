package main

import (
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

type telenetClient struct {
	connection net.Conn
	timeout    time.Duration

	in  io.ReadCloser
	out io.Writer

	protocol string
	address  string
}

func readWrite(reader io.Reader, writer io.Writer) error {
	message, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	_, err = writer.Write(message)
	return err
}

func (client *telenetClient) Receive() error {
	return readWrite(client.connection, client.out)
}

func (client *telenetClient) Send() error {
	return readWrite(client.in, client.connection)
}

func (client *telenetClient) Close() error {
	return client.connection.Close()
}

func (client *telenetClient) Connect() error {
	conn, err := net.DialTimeout(client.protocol, client.address, client.timeout)
	if err != nil {
		return err
	}
	client.connection = conn
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telenetClient{
		protocol: "tcp",
		address:  address,
		timeout:  timeout,
		in:       in,
		out:      out,
	}
}
