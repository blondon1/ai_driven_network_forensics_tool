package main

import (
	"github.com/blondon1/ai_driven_network_forensics_tool/src/analysis"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/data_ingestion"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/preprocessing"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/real_time_analysis"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/reporting"
)

func main() {
	packets := make(chan gopacket.Packet)

	go func() {
		data_ingestion.CapturePackets("eth0", packets)
	}()

	for packet := range packets {
		preprocessing.PreprocessPacket(packet)
		analysis.AnalyzePacket(packet)
		real_time_analysis.AnalyzeInRealTime(packet)
		reporting.GenerateReport(packet)
	}
}
