package main

import (
	"fmt"
	"github.com/google/gopacket/pcap"
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Printf("Error finding devices: %v\n", err)
		return
	}

	fmt.Println("Devices found:")
	for _, device := range devices {
		fmt.Printf("- Name: %s, Description: %s\n", device.Name, device.Description)
	}
}
