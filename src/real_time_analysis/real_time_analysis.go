package real_time_analysis

import (
	"ai_driven_network_forensics_tool/src/ui"
	"fmt"
	"github.com/google/gopacket"
	"log"
	"os/exec"
)

// Analyzes packets in real-time and triggers alerts for suspicious activity
func AnalyzeInRealTime(packet gopacket.Packet) {
	packetSize := len(packet.Data())

	// Record traffic data for visualization
	ui.RecordPacketCount()

	if packetSize > 1000 { // Example threshold for triggering an alert
		TriggerAlert(packetSize)
		ExecuteResponseActions("HighPacketSizeDetected", fmt.Sprintf("Packet size: %d bytes", packetSize))
	}

	fmt.Println("Real-time analysis of packet:", packet)
}

// TriggerAlert sends an alert if suspicious activity is detected
func TriggerAlert(packetSize int) {
	message := fmt.Sprintf("ALERT: Suspicious packet size detected: %d bytes", packetSize)
	log.Println(message)
	ui.BroadcastAlert(message)
}

// ExecuteResponseActions executes predefined scripts or actions when an anomaly is detected
func ExecuteResponseActions(actionType, details string) {
	log.Printf("Executing response action: %s - Details: %s\n", actionType, details)
	// Example: Blocking an IP or notifying an admin
	switch actionType {
	case "HighPacketSizeDetected":
		// Example script execution: Notify admin
		cmd := exec.Command("sh", "-c", "echo 'High packet size detected: "+details+"' | mail -s 'Network Alert' admin@example.com")
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to execute response action: %v\n", err)
		}
	}
}
