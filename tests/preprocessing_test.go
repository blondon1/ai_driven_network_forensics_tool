package tests

import (
	"github.com/blondon1/ai_driven_network_forensics_tool/src/preprocessing"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"testing"
)

func TestPreprocessPacket(t *testing.T) {
	// This is a placeholder test function
	t.Log("Running TestPreprocessPacket...")

	packet := gopacket.NewPacket([]byte{}, layers.LayerTypeEthernet, gopacket.Default)
	preprocessing.PreprocessPacket(packet)
}
