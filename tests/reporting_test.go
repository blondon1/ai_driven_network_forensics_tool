package tests

import (
	"github.com/blondon1/ai_driven_network_forensics_tool/src/reporting"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"testing"
)

func TestGenerateReport(t *testing.T) {
	t.Log("Running TestGenerateReport...")

	packetData := make([]byte, 100)
	packet := gopacket.NewPacket(packetData, layers.LayerTypeEthernet, gopacket.Default)

	reporting.GenerateReport(packet)

	// Check if the report was generated successfully
	// You could further expand this test by checking if the file exists, etc.
}
