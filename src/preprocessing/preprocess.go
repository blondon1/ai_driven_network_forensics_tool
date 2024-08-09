package preprocessing

import (
	"fmt"
	"github.com/google/gopacket"
)

func PreprocessPacket(packet gopacket.Packet) {
	fmt.Println("Packet Length:", len(packet.Data()))
}
