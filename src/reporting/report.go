package reporting

import (
	"fmt"
	"github.com/google/gopacket"
	"os"
	"time"
)

// Generates a detailed forensic report for the given packet
func GenerateReport(packet gopacket.Packet) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	report := fmt.Sprintf("Timestamp: %s\nPacket Length: %d\nPacket Data: %x\n\n",
		timestamp, len(packet.Data()), packet.Data())

	// Save the report to a file
	fileName := fmt.Sprintf("data/logs/report_%s.txt", time.Now().Format("20060102_150405"))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to create report file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(report)
	if err != nil {
		fmt.Println("Failed to write to report file:", err)
	}
}
