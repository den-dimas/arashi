package main

import (
	"log"
	"time"

	"github.com/simonvetter/modbus"
	"gitlab.com/rt.ece.ntust/rt-lab/arashi/internal/utils"
	device "gitlab.com/rt.ece.ntust/rt-lab/arashi/pkg/device/twintex-2000v"
)

func main() {
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:      "rtu://COM3",
		Speed:    9600,
		DataBits: 8,
		StopBits: 1,
		Parity:   modbus.PARITY_NONE,
		Timeout:  5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = client.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	client.SetUnitId(1)

	readState := func(label string) {
		time.Sleep(150 * time.Millisecond)
		res, err := client.ReadRegisters(device.REG_VOLT_SET, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			log.Fatalf("Read error: %v", err)
		}
		v := utils.DecodeFloat(res[0], res[1])
		i := utils.DecodeFloat(res[2], res[3])
		log.Printf("[%s] Status -> Voltage: %.2fV, Current: %.4fA", label, v, i)
	}

	err = client.WriteRegisters(device.REG_CURR_SET, utils.EncodeFloat(100.0))
	if err != nil {
		log.Fatal(err)
	}
	readState("Step 1: Set Current 100mA")

	err = client.WriteRegister(device.REG_CONTROL, device.CMD_CURR_ENABLE)
	if err != nil {
		log.Fatal(err)
	}
	readState("Step 2: Enable Current")

	err = client.WriteRegisters(device.REG_VOLT_SET, utils.EncodeFloat(0.0))
	if err != nil {
		log.Fatal(err)
	}
	readState("Step 3: Set Voltage 0V")

	err = client.WriteRegister(device.REG_CONTROL, device.CMD_VOLT_ENABLE)
	if err != nil {
		log.Fatal(err)
	}
	readState("Step 4: Enable Voltage")

	readState("Before Step 5: Output ON")
	err = client.WriteRegister(device.REG_CONTROL, device.CMD_OUTPUT_ON)
	if err != nil {
		log.Fatal(err)
	}
	readState("After Step 5: Output ON")

	err = client.WriteRegisters(device.REG_VOLT_SET, utils.EncodeFloat(10.0))
	if err != nil {
		log.Fatal(err)
	}
	readState("Step 6: Set Voltage 10V")

	err = client.WriteRegisters(device.REG_RAMP_TIME, utils.EncodeFloat(10.0))
	if err != nil {
		log.Fatal(err)
	}
	readState("Step 7: Set Ramp 10s")

	err = client.WriteRegister(device.REG_CONTROL, device.CMD_RAMP_ENABLE)
	if err != nil {
		log.Fatal(err)
	}
	readState("Step 8: Enable Ramp")

	readState("Before Step 9: Final ON")
	err = client.WriteRegister(device.REG_CONTROL, device.CMD_OUTPUT_ON)
	if err != nil {
		log.Fatal(err)
	}
	readState("After Step 9: Final ON")

	for j := 0; j < 12; j++ {
		time.Sleep(1 * time.Second)
		readState("Monitoring Ramp Progress")
	}
}
