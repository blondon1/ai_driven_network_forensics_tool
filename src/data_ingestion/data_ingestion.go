package data_ingestion

import (
	"database/sql"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

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
			storePacketMetadata(packet) // Store metadata
			packets <- packet
		}
	}
}

// StorePacketMetadata saves packet metadata to the SQLite database
func storePacketMetadata(packet gopacket.Packet) {
	db, err := sql.Open("sqlite3", "data/packets.db")
	if err != nil {
		log.Println("Failed to open database:", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS packets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME,
		source_ip TEXT,
		destination_ip TEXT,
		protocol TEXT,
		size INTEGER
	)`)
	if err != nil {
		log.Println("Failed to create table:", err)
		return
	}

	ipLayer := packet.Layer(gopacket.LayerTypeIPv4)
	if ipLayer == nil {
		log.Println("No IP layer found in packet")
		return
	}

	ip, _ := ipLayer.(*gopacket.LayerTypeIPv4)
	protocol := ip.Protocol.String()
	size := len(packet.Data())

	_, err = db.Exec(`INSERT INTO packets (timestamp, source_ip, destination_ip, protocol, size)
		VALUES (?, ?, ?, ?, ?)`,
		time.Now(),
		ip.SrcIP.String(),
		ip.DstIP.String(),
		protocol,
		size,
	)
	if err != nil {
		log.Println("Failed to insert packet data:", err)
	}
}
