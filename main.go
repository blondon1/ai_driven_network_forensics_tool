package main

import (
	"fmt"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/analysis"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/data_ingestion"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/preprocessing"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/real_time_analysis"
	"github.com/blondon1/ai_driven_network_forensics_tool/src/reporting"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type SystemConfig struct {
	NetworkInterface string `yaml:"network_interface"`
	LogFilePath      string `yaml:"log_file_path"`
}

func loadConfig() (SystemConfig, error) {
	var config SystemConfig
	file, err := os.Open("config/system_config.yaml")
	if err != nil {
		return config, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
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
		log.Fatal("Failed to load configuration:", err)
	}

	logFile, err := setupLogging(config.LogFilePath)
	if err != nil {
		log.Fatal("Failed to set up logging:", err)
	}
	defer logFile.Close()

	log.Println("Starting AI-Driven Network Forensics Tool")

	packets := make(chan gopacket.Packet)

	go func() {
		data_ingestion.CapturePackets(config.NetworkInterface, packets)
	}()

	for packet := range packets {
		preprocessing.PreprocessPacket(packet)
		analysis.AnalyzePacket(packet)
		real_time_analysis.AnalyzeInRealTime(packet)
		reporting.GenerateReport(packet)
	}

	log.Println("Shutting down AI-Driven Network Forensics Tool")
}
