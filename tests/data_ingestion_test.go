package tests

import (
	"github.com/blondon1/ai_driven_network_forensics_tool/src/analysis"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"testing"
)

func TestAnalyzePacket(t *testing.T) {
	t.Log("Running TestAnalyzePacket...")

	packetData := make([]byte, 100)
	packet := gopacket.NewPacket(packetData, layers.LayerTypeEthernet, gopacket.Default)

	for i := 0; i < 15; i++ {
		analysis.AnalyzePacket(packet)
	}

	// Test should check for an anomaly after sufficient identical packets
	// Since the current logic flags anomalies after 10 identical packets,
	// This test should trigger the anomaly detection.
}
