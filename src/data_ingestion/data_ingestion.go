package data_ingestion

import (
    "database/sql"
    "log"
    "time"
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
    _ "github.com/mattn/go-sqlite3"
)

type FilterConfig struct {
    BPF       string  // Added BPF as a string field for Berkeley Packet Filter expressions
    IPAddress string
    Port      int
    Protocol  string
}

// Modify the CapturePackets function to use BPF
func CapturePackets(interfaceName string, packets chan<- gopacket.Packet, filterConfig FilterConfig) {
    handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
    if err != nil {
        log.Fatal(err)
    }
    defer handle.Close()

    if filterConfig.BPF != "" {
        if err := handle.SetBPFFilter(filterConfig.BPF); err != nil {
            log.Fatal("Error setting BPF filter:", err)
        }
    }

    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
        log.Println("Packet captured:", packet) // Debug log for each captured packet
        packets <- packet
    }
}


func SavePacketToDatabase(db *sql.DB, packet gopacket.Packet) {
    ipLayer := packet.Layer(layers.LayerTypeIPv4)
    tcpLayer := packet.Layer(layers.LayerTypeTCP)

    if ipLayer != nil && tcpLayer != nil {
        ip, _ := ipLayer.(*layers.IPv4)
        tcp, _ := tcpLayer.(*layers.TCP)

        timestamp := time.Now().Format("2006-01-02 15:04:05")
        query := `
            INSERT INTO packets (timestamp, source_ip, destination_ip, protocol, source_port, destination_port, packet_size)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        `

        _, err := db.Exec(query, timestamp, ip.SrcIP, ip.DstIP, "TCP", tcp.SrcPort, tcp.DstPort, len(packet.Data()))
        if err != nil {
            log.Printf("Failed to save packet to database: %v", err)
        }
    }
}

