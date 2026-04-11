package device

import (
	"math/rand"
	"sync"
	"time"

	"github.com/simonvetter/modbus"
	"gitlab.com/rt.ece.ntust/rt-lab/arashi/internal/utils"
	device "gitlab.com/rt.ece.ntust/rt-lab/arashi/pkg/device/protocol"
)

type Twintex2000V struct {
	lock           sync.RWMutex
	targetVoltage  float32
	targetCurrent  float32
	targetRampTime float32
	outputEnabled  bool

	voltage float32
	current float32

	rampStartVoltage float32
	rampStartTime    time.Time
	isRamping        bool

	pendingVoltage  float32
	pendingCurrent  float32
	pendingRampTime float32
}

const (
	REG_VOLT_MEAS = 0x0B00
	REG_CURR_MEAS = 0x0B02

	REG_CONTROL   = 0x0A00
	REG_VOLT_SET  = 0x0A05
	REG_CURR_SET  = 0x0A07
	REG_RAMP_TIME = 0x0A09

	CMD_VOLT_ENABLE = 1
	CMD_CURR_ENABLE = 2
	CMD_RAMP_ENABLE = 3
	CMD_OUTPUT_ON   = 6
)

func NewTwintex2000V(voltage, current float32) *Twintex2000V {
	return &Twintex2000V{
		targetVoltage: voltage,
		targetCurrent: current,
	}
}

func (t *Twintex2000V) Start(config *modbus.ServerConfiguration, tick time.Duration) error {
	mb := device.NewModbusDevice(t, t.Process)

	err := mb.Start(config, tick)
	if err != nil {
		return err
	}

	return nil
}

func (t *Twintex2000V) Process() {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.outputEnabled {
		t.voltage = 0.0
		t.current = 0.0
		t.isRamping = false
		return
	}

	if t.isRamping && t.targetRampTime > 0 {
		elapsed := time.Since(t.rampStartTime).Seconds()
		progress := float32(elapsed) / t.targetRampTime

		if progress >= 1.0 {
			t.voltage = t.targetVoltage
			t.isRamping = false
		} else {
			t.voltage = t.rampStartVoltage + ((t.targetVoltage - t.rampStartVoltage) * progress)
		}
	} else {
		t.voltage = t.targetVoltage
	}

	t.voltage += (rand.Float32() * 0.2)
	t.current = t.targetCurrent + (rand.Float32() * 0.05)
}

func (t *Twintex2000V) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if req.IsWrite {
		return t.handleWrites(req)
	}

	return t.handleReads(req)
}

func (t *Twintex2000V) handleWrites(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	switch req.Addr {
	case REG_CONTROL:
		command := req.Args[0]
		switch command {
		case CMD_VOLT_ENABLE:
			t.targetVoltage = t.pendingVoltage
			t.isRamping = false
		case CMD_CURR_ENABLE:
			t.targetCurrent = t.pendingCurrent
		case CMD_RAMP_ENABLE:
			t.targetVoltage = t.pendingVoltage
			t.targetCurrent = t.pendingCurrent
			t.targetRampTime = t.pendingRampTime

			t.rampStartVoltage = t.voltage
			t.rampStartTime = time.Now()
			t.isRamping = true
		case CMD_OUTPUT_ON:
			t.outputEnabled = true
		}

	case REG_VOLT_SET:
		t.pendingVoltage = utils.DecodeFloat(req.Args[0], req.Args[1])
	case REG_CURR_SET:
		t.pendingCurrent = utils.DecodeFloat(req.Args[0], req.Args[1])
	case REG_RAMP_TIME:
		t.pendingRampTime = utils.DecodeFloat(req.Args[0], req.Args[1])
	}

	return nil, nil
}

func (t *Twintex2000V) handleReads(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	for i := uint16(0); i < req.Quantity; i++ {
		addr := req.Addr + i

		switch addr {
		case REG_VOLT_SET:
			res = append(res, utils.EncodeFloat(t.voltage)...)
			i++
		case REG_CURR_SET:
			res = append(res, utils.EncodeFloat(t.current)...)
			i++
		default:
			res = append(res, 0)
		}
	}

	return res, nil
}

func (t *Twintex2000V) HandleCoils(req *modbus.CoilsRequest) (res []bool, err error) {
	err = modbus.ErrIllegalFunction

	return
}

func (t *Twintex2000V) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error) {
	err = modbus.ErrIllegalFunction

	return
}

func (t *Twintex2000V) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	err = modbus.ErrIllegalFunction
	return
}
