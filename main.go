package main

import (
	"github.com/blondon1/ai_driven_network_forensics_tool/src/analysis"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/data_ingestion"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/preprocessing"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/real_time_analysis"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/reporting"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/ui"
	"github.com/google/gopacket"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type SystemConfig struct {
	NetworkInterface string                      `yaml:"network_interface"`
	LogFilePath      string                      `yaml:"log_file_path"`
	Filter           data_ingestion.FilterConfig `yaml:"filter"`
	Authentication   struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"authentication"`
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
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return logFile, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ui.username = config.Authentication.Username
	ui.password = config.Authentication.Password

	logFile, err := setupLogging(config.LogFilePath)
	if err != nil {
		log.Fatalf("Failed to set up logging: %v", err)
	}
	defer logFile.Close()

	log.Println("Starting AI-Driven Network Forensics Tool")

	packets := make(chan gopacket.Packet)

	go func() {
		data_ingestion.CapturePackets(config.NetworkInterface, packets, config.Filter)
	}()

	go func() {
		ui.StartServer()
	}()

	for packet := range packets {
		preprocessing.PreprocessPacket(packet)
		analysis.AnalyzePacket(packet)
		real_time_analysis.AnalyzeInRealTime(packet)
		reporting.GenerateReport(packet)
	}

	log.Println("Shutting down AI-Driven Network Forensics Tool")
}
