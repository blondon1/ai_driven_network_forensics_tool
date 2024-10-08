package main

import (
	"github.com/blondon1/ai_driven_network_forensics_tool/src/analysis"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/data_ingestion"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/preprocessing"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/real_time_analysis"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/reporting"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/ui"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers" // Import for LayerTypeTCP and LayerTypeUDP
	"github.com/google/gopacket/pcap"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

type SystemConfig struct {
	NetworkInterface string `yaml:"network_interface"`
	LogFilePath      string `yaml:"log_file_path"`
}

func loadConfig() (SystemConfig, error) {
	var config SystemConfig
	data, err := os.ReadFile("config/system_config.yaml")
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func setupLogging(logFilePath string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(logFilePath), os.ModePerm)
	if err != nil {
		return nil, err
	}
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return logFile, nil
}

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Devices found:")
	for _, device := range devices {
		log.Printf("Name: %s, Description: %s\n", device.Name, device.Description)
	}

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logFile, err := setupLogging(config.LogFilePath)
	if err != nil {
		log.Fatalf("Failed to set up logging: %v", err)
	}
	defer logFile.Close()

	log.Println("Starting AI-Driven Network Forensics Tool")

	packets := make(chan gopacket.Packet)

	filterConfig := data_ingestion.FilterConfig{
		BPF:       "",
		IPAddress: "",
		Port:      0,
		Protocol:  "",
	}

	go func() {
		data_ingestion.CapturePackets(config.NetworkInterface, packets, filterConfig)
	}()

	go func() {
		ui.StartServer()
	}()

	for packet := range packets {
		log.Println("Processing packet")
		preprocessing.PreprocessPacket(packet)
		analysis.AnalyzePacket(packet)

		protocol := "Other"
		if packet.TransportLayer() != nil {
			switch packet.TransportLayer().LayerType() {
			case layers.LayerTypeTCP:
				protocol = "TCP"
			case layers.LayerTypeUDP:
				protocol = "UDP"
			default:
				if packet.Layer(layers.LayerTypeICMPv4) != nil || packet.Layer(layers.LayerTypeICMPv6) != nil {
					protocol = "ICMP"
				}
			}
		}
		log.Printf("Packet protocol: %s", protocol)
		ui.RecordPacketCount(protocol)
		real_time_analysis.AnalyzeInRealTime(packet)
		reporting.GenerateReport(packet)
	}

	log.Println("Shutting down AI-Driven Network Forensics Tool")
}
