// Package protocol contains all of the supported protocol for the IoT devices.
package protocol

import (
	"bufio"
	"net"
	"strings"
	"time"
)

type SCPIClient struct {
	address string
	conn    net.Conn
	reader  *bufio.Reader
}

func NewSCPIClient(address string, timeout time.Duration) (*SCPIClient, error) {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, err
	}

	return &SCPIClient{
		address: address,
		conn:    conn,
		reader:  bufio.NewReader(conn),
	}, nil
}

func (c *SCPIClient) Close() error {
	return c.conn.Close()
}

func (c *SCPIClient) Write(cmd string) error {
	cmd = strings.TrimSpace(cmd) + "\r\n"
	_, err := c.conn.Write([]byte(cmd))
	return err
}

func (c *SCPIClient) Query(cmd string) (string, error) {
	err := c.Write(cmd)
	if err != nil {
		return "", err
	}

	response, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}
