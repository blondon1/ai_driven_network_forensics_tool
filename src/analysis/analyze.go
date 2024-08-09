package analysis

import (
	"fmt"
	"github.com/google/gopacket"
)

// Basic anomaly detection: counts packet sizes and flags unusual sizes
var packetSizeCounts = make(map[int]int)

func AnalyzePacket(packet gopacket.Packet) {
	packetSize := len(packet.Data())
	packetSizeCounts[packetSize]++

	// Simple anomaly detection: flag packet sizes with unusually high counts
	if packetSizeCounts[packetSize] > 10 {
		fmt.Printf("Anomaly detected: packet size %d has appeared %d times\n", packetSize, packetSizeCounts[packetSize])
	}
}
