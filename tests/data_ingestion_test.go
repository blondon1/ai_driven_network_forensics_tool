package tests

import (
	"github.com/blondon1/ai_driven_network_forensics_tool/src/data_ingestion"
	"github.com/google/gopacket"
	"testing"
)

func TestCapturePackets(t *testing.T) {
	t.Log("Running TestCapturePackets...")

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CapturePackets panicked: %v", r)
			}
		}()
		data_ingestion.CapturePackets("eth0", make(chan gopacket.Packet))
	}()
}
