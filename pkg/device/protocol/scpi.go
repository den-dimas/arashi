// Package device is for
package device

import (
	"bufio"
	"net"
	"strings"
)

type SCPIDevice struct {
	address string
	isOn    bool
	isTrip  bool

	ProcessCommand func(string) string
}

func NewSCPIDevice(address string, ProcessCommand func(string) string) *SCPIDevice {
	return &SCPIDevice{
		address:        address,
		isOn:           false,
		isTrip:         false,
		ProcessCommand: ProcessCommand,
	}
}

func (d *SCPIDevice) Start() error {
	ln, err := net.Listen("tcp", d.address)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go d.handleConnection(conn)
	}
}

func (d *SCPIDevice) handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		response := d.ProcessCommand(input)
		if response != "" {
			conn.Write([]byte(response + "\r\n"))
		}
	}
}
