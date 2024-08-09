package tests

import (
	"github.com/blondon1/ai_driven_network_forensics_tool/src/real_time_analysis"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"testing"
)

func TestAnalyzeInRealTime(t *testing.T) {
	t.Log("Running TestAnalyzeInRealTime...")

	packetData := make([]byte, 1500) // This should trigger an alert
	packet := gopacket.NewPacket(packetData, layers.LayerTypeEthernet, gopacket.Default)

	real_time_analysis.AnalyzeInRealTime(packet)

	// You could further expand this test by capturing log output, etc.
}
