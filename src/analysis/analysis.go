package analysis

import (
	"fmt"
	"github.com/google/gopacket"
	"log"
	"sync"
	"time"
)

// Packet size statistics
type PacketSizeStats struct {
	Count int
	Total int
}

var (
	packetSizeCounts = make(map[int]int)
	packetStats      = make(map[time.Time]PacketSizeStats)
	mu               sync.Mutex
)

// AnalyzePacket performs anomaly detection on each packet
func AnalyzePacket(packet gopacket.Packet) {
	packetSize := len(packet.Data())
	mu.Lock()
	packetSizeCounts[packetSize]++
	currentTime := time.Now().Truncate(time.Minute)

	// Update stats for the current time window
	if stats, exists := packetStats[currentTime]; exists {
		stats.Count++
		stats.Total += packetSize
		packetStats[currentTime] = stats
	} else {
		packetStats[currentTime] = PacketSizeStats{Count: 1, Total: packetSize}
	}
	mu.Unlock()

	// Detect anomalies
	if packetSizeCounts[packetSize] > 10 {
		fmt.Printf("Anomaly detected: packet size %d has appeared %d times\n", packetSize, packetSizeCounts[packetSize])
	}

	if isStatisticalAnomaly(packetSize, currentTime) {
		fmt.Printf("Statistical anomaly detected: packet size %d at %v\n", packetSize, currentTime)
	}

	if isTimeAnomaly(currentTime) {
		fmt.Printf("Time-based anomaly detected at %v\n", currentTime)
	}
}

// isStatisticalAnomaly checks for deviations from the average packet size
func isStatisticalAnomaly(packetSize int, currentTime time.Time) bool {
	mu.Lock()
	defer mu.Unlock()
	stats := packetStats[currentTime]
	averageSize := stats.Total / stats.Count
	return packetSize > averageSize*2 || packetSize < averageSize/2
}

// isTimeAnomaly detects spikes in traffic volume within short time windows
func isTimeAnomaly(currentTime time.Time) bool {
	mu.Lock()
	defer mu.Unlock()
	previousTime := currentTime.Add(-time.Minute)
	if previousStats, exists := packetStats[previousTime]; exists {
		currentStats := packetStats[currentTime]
		return currentStats.Count > previousStats.Count*2
	}
	return false
}
