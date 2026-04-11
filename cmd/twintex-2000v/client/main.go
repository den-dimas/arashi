package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/simonvetter/modbus"
	"gitlab.com/rt.ece.ntust/rt-lab/arashi/internal/utils"
	device "gitlab.com/rt.ece.ntust/rt-lab/arashi/pkg/device/twintex-2000v"
)

func main() {
	portPtr := flag.String("port", "COM6", "USB Port connected to the Twintex2000V")

	flag.Parse()

	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:      fmt.Sprintf("rtu://%s", *portPtr),
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

	err = client.WriteCoil(0x0500, true)
	if err != nil {
		log.Printf("Failed to set remote control on: %v", err)
	}
	log.Printf("Set remote control on")

	err = client.WriteRegisters(device.REG_CURR_SET, utils.EncodeFloat(100.0))
	if err != nil {
		log.Fatal(err)
	}
	readData(device.REG_CURR_SET, "CURRENT SET REGISTER", client)

	err = client.WriteRegisters(device.REG_CONTROL, []uint16{device.CMD_CURR_ENABLE})
	if err != nil {
		log.Fatal(err)
	}
	readData(device.REG_CURR_MEAS, "ACTUAL CURRENT VALUE", client)

	err = client.WriteRegisters(device.REG_VOLT_SET, utils.EncodeFloat(0.0))
	if err != nil {
		log.Fatal(err)
	}
	readData(device.REG_VOLT_SET, "VOLTAGE SET REGISTER", client)

	err = client.WriteRegisters(device.REG_CONTROL, []uint16{device.CMD_VOLT_ENABLE})
	if err != nil {
		log.Fatal(err)
	}
	readData(device.REG_VOLT_MEAS, "ACTUAL VOLTAGE VALUE", client)

	err = client.WriteRegisters(device.REG_CONTROL, []uint16{device.CMD_OUTPUT_ON})
	if err != nil {
		log.Fatal(err)
	}
	readData(device.REG_VOLT_MEAS, "ACTUAL VOLTAGE VALUE", client)
	readData(device.REG_CURR_MEAS, "ACTUAL CURRENT VALUE", client)

	err = client.WriteRegisters(device.REG_VOLT_SET, utils.EncodeFloat(10.0))
	if err != nil {
		log.Fatal(err)
	}
	readData(device.REG_VOLT_SET, "VOLTAGE SET REGISTER", client)

	err = client.WriteRegisters(device.REG_RAMP_TIME, utils.EncodeFloat(10.0))
	if err != nil {
		log.Fatal(err)
	}
	readData(device.REG_RAMP_TIME, "RAMP TIME SET REGISTER", client)

	err = client.WriteRegisters(device.REG_CONTROL, []uint16{device.CMD_RAMP_ENABLE})
	if err != nil {
		log.Fatal(err)
	}

	err = client.WriteRegisters(device.REG_CONTROL, []uint16{device.CMD_OUTPUT_ON})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("\nFinal Command On\n")

	for j := 0; j < 12; j++ {
		time.Sleep(1 * time.Second)
		readData(device.REG_CURR_MEAS, "ACTUAL CURRENT VALUE", client)
		readData(device.REG_VOLT_MEAS, "ACTUAL VOLTAGE VALUE", client)
	}
}

func readData(address uint16, label string, client *modbus.ModbusClient) float32 {
	time.Sleep(150 * time.Millisecond)

	var dataType string
	var m string

	if address == device.REG_VOLT_MEAS || address == device.REG_VOLT_SET {
		dataType = "Voltage"
		m = "V"
	} else if address == device.REG_CURR_MEAS || address == device.REG_CURR_SET {
		dataType = "Current"
		m = "mA"
	} else {
		dataType = "Ramp Time"
		m = "s"
	}

	dataRaw, err := client.ReadRegisters(address, 2, modbus.HOLDING_REGISTER)
	if err != nil {
		log.Fatalf("Read %s error: %v", dataType, err)
	}

	d := utils.DecodeFloat(dataRaw[0], dataRaw[1])

	log.Printf("[%s] %s: %.2f%s", label, dataType, d, m)
	return d
}
