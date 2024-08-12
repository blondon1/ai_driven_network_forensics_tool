package ai

import (
    "bytes"
    "fmt"
    "log"
    "os/exec"
    "strings"
)

func IsAnomalous(packetLength int) bool {
    cmd := exec.Command("python", "src/ai/analyze_packet.py", "--packet-length", fmt.Sprint(packetLength))
    
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr

    err := cmd.Run()
    if err != nil {
        log.Printf("Failed to run Python script: %s\n", stderr.String())
        return false
    }

    output := out.String()
    log.Printf("Python script output: %s", output)

    // Check if the output contains "Anomaly"
    return strings.Contains(output, "Anomaly")
}
