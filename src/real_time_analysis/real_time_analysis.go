package real_time_analysis

import (
    "github.com/blondon1/ai_driven_network_forensics_tool/src/ui"
    "github.com/google/gopacket"
    "log"
)

// AnalyzeInRealTime analyzes packets in real-time and triggers alerts
func AnalyzeInRealTime(packet gopacket.Packet) {
    packetSize := len(packet.Data())

    if packetSize > 1000 { // Example threshold for triggering an alert
        alertMessage := "Suspicious packet detected with size: " + string(packetSize) + " bytes"
        log.Println(alertMessage)
        ui.SendAlert(alertMessage) // Send alert to be broadcasted to the dashboard
    }
}
