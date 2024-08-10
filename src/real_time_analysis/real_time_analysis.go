package real_time_analysis

import (
	"ai_driven_network_forensics_tool/src/ui"
	"fmt"
	"github.com/google/gopacket"
	"log"
)

// Analyzes packets in real-time and triggers alerts for suspicious activity
func AnalyzeInRealTime(packet gopacket.Packet) {
	packetSize := len(packet.Data())

	// Record traffic data for visualization
	ui.RecordPacketCount()

	if packetSize > 1000 { // Example threshold for triggering an alert
		TriggerAlert(packetSize)
	}

	fmt.Println("Real-time analysis of packet:", packet)
}

// TriggerAlert sends an alert if suspicious activity is detected
func TriggerAlert(packetSize int) {
	message := fmt.Sprintf("ALERT: Suspicious packet size detected: %d bytes", packetSize)
	log.Println(message)
	ui.BroadcastAlert(message)
}
