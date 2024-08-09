package data_ingestion

import (
    "log"
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
)

func CapturePackets(interfaceName string, packets chan<- gopacket.Packet) {
    handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
    if err != nil {
        log.Fatal(err)
    }
    defer handle.Close()

    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
        packets <- packet
    }
}

