// Package device
package device

import (
	"fmt"
	"strings"

	device "gitlab.com/rt.ece.ntust/rt-lab/arashi/pkg/device/protocol"
)

type PFR100M struct {
	voltage float64
	current float64
}

const (
	MEASURE    = ":MEAS:ALL:DC?"
	CURR_LIMIT = ":SOUR:CURR:LEV:IMM:AMPL"
)

func NewPFRM100M() *PFR100M {
	return &PFR100M{}
}

func (p *PFR100M) Start(address string) error {
	device := device.NewSCPIDevice(address, p.processCommand)

	err := device.Start()
	if err != nil {
		return err
	}

	return nil
}

func (p *PFR100M) processCommand(cmd string) string {
	upperCmd := strings.ToUpper(cmd)

	switch {
	case upperCmd == MEASURE:
		vStr := fmt.Sprintf("%.3f", p.voltage)
		iStr := fmt.Sprintf("%.4f", p.current)
		return "-" + vStr + ",+" + iStr

	case strings.HasPrefix(upperCmd, CURR_LIMIT):
		parts := strings.Fields(cmd)
		if len(parts) > 1 {
			fmt.Sscanf(parts[1], "%f", &p.voltage)
		}
		return ""
	default:
		return ""
	}
}
