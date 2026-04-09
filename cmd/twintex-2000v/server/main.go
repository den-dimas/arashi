package main

import (
	"log"
	"os"
	"time"

	"github.com/simonvetter/modbus"
	device "gitlab.com/rt.ece.ntust/rt-lab/arashi/pkg/device/twintex-2000v"
)

func main() {
	device := device.NewTwintex2000V(220.0, 5.0)

	address := "0.0.0.0:4832"

	log.Printf("Starting Twintex2000V Modbus Server on %s...", address)

	err := device.Start(&modbus.ServerConfiguration{
		URL:        address,
		Timeout:    10 * time.Second,
		MaxClients: 5,
	}, 100*time.Millisecond)
	if err != nil {
		log.Printf("Error starting Twintex2000V Server: %v", err)
		os.Exit(1)
	}

	log.Printf("Twintex2000V Server has started at %s.", address)
}
