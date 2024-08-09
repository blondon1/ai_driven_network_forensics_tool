package main

import (
	"github.com/blondon1/ai_driven_network_forensics_tool/src/data_ingestion"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/preprocessing"
)

func main() {
	packets := make(chan gopacket.Packet)

	go func() {
		data_ingestion.CapturePackets("eth0", packets)
	}()

	for packet := range packets {
		preprocessing.PreprocessPacket(packet)
	}
}
