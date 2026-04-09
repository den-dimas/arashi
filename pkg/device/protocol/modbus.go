package device

import (
	"time"

	"github.com/simonvetter/modbus"
)

type ModbusDevice struct {
	Handler modbus.RequestHandler
	Process func()
}

func NewModbusDevice(Handler modbus.RequestHandler, Process func()) *ModbusDevice {
	return &ModbusDevice{
		Handler: Handler,
		Process: Process,
	}
}

func (m *ModbusDevice) Start(config *modbus.ServerConfiguration, tick time.Duration) error {
	server, err := modbus.NewServer(&modbus.ServerConfiguration{
		URL:        "tcp://" + config.URL,
		Timeout:    config.Timeout,
		MaxClients: config.MaxClients,
	}, m.Handler)
	if err != nil {
		return err
	}

	err = server.Start()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(tick)
	for range ticker.C {
		if m.Process != nil {
			m.Process()
		}
	}

	return nil
}
