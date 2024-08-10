package data_ingestion

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"strings"
)

// CapturePackets captures packets from the network interface and applies filtering
func CapturePackets(interfaceName string, packets chan<- gopacket.Packet, filterConfig FilterConfig) {
	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Apply BPF filter if specified
	if filterConfig.BPF != "" {
		err := handle.SetBPFFilter(filterConfig.BPF)
		if err != nil {
			log.Fatalf("Failed to set BPF filter: %v", err)
		}
		log.Printf("BPF filter applied: %s", filterConfig.BPF)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if filterPacket(packet, filterConfig) {
			packets <- packet
		}
	}
}

// FilterConfig holds the filtering criteria
type FilterConfig struct {
	SourceIP      string
	DestinationIP string
	Protocol      string
	Port          string
	BPF           string // Berkeley Packet Filter string for advanced filtering
}

// filterPacket filters packets based on the criteria in FilterConfig
func filterPacket(packet gopacket.Packet, config FilterConfig) bool {
	if config.SourceIP != "" {
		if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)
			if !strings.Contains(ip.SrcIP.String(), config.SourceIP) {
				return false
			}
		}
	}

	if config.DestinationIP != "" {
		if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)
			if !strings.Contains(ip.DstIP.String(), config.DestinationIP) {
				return false
			}
		}
	}

	if config.Protocol != "" {
		if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)
			if !strings.EqualFold(ip.Protocol.String(), config.Protocol) {
				return false
			}
		}
	}

	if config.Port != "" {
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			if tcp.SrcPort.String() != config.Port && tcp.DstPort.String() != config.Port {
				return false
			}
		}
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			udp, _ := udpLayer.(*layers.UDP)
			if udp.SrcPort.String() != config.Port && udp.DstPort.String() != config.Port {
				return false
			}
		}
	}

	return true
}
