package ui

import (
    "github.com/go-echarts/go-echarts/v2/charts"
    "github.com/go-echarts/go-echarts/v2/components"
    "github.com/go-echarts/go-echarts/v2/opts"
    "github.com/gorilla/websocket"
    "log"
    "net/http"
    "sync"
    "time"
)

var (
    alertsMutex         sync.Mutex
    alertsConn          []*websocket.Conn
    packetCountMutex    sync.Mutex
    packetCountPerMinute = make(map[string]int)
    protocolCounts      = map[string]int{
        "TCP":  0,
        "UDP":  0,
        "ICMP": 0,
        "Other": 0,
    }
    upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        CheckOrigin:     func(r *http.Request) bool { return true },
    }
)

func StartServer() {
    http.HandleFunc("/", dashboardHandler)
    http.HandleFunc("/ws", websocketHandler)
    log.Println("Starting web server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    page := components.NewPage()
    page.AddCharts(generateTrafficLineChart(), generateProtocolPieChart(), generateAnomalyChart())
    page.Render(w)
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Failed to set up WebSocket:", err)
        return
    }
    defer conn.Close()
    alertsMutex.Lock()
    alertsConn = append(alertsConn, conn)
    alertsMutex.Unlock()

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            break
        }
    }
}

func generateTrafficLineChart() *charts.Line {
    line := charts.NewLine()
    line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
        Title: "Network Traffic",
    }))

    packetCountMutex.Lock()
    defer packetCountMutex.Unlock()

    var times []string
    var counts []opts.LineData
    for timeStr, count := range packetCountPerMinute {
        times = append(times, timeStr)
        counts = append(counts, opts.LineData{Value: count})
    }

    line.SetXAxis(times).
        AddSeries("Packets per Minute", counts)

    return line
}

func generateProtocolPieChart() *charts.Pie {
    pie := charts.NewPie()
    pie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
        Title: "Protocol Distribution",
    }))

    pie.AddSeries("Protocol", []opts.PieData{
        {Name: "TCP", Value: protocolCounts["TCP"]},
        {Name: "UDP", Value: protocolCounts["UDP"]},
        {Name: "ICMP", Value: protocolCounts["ICMP"]},
        {Name: "Other", Value: protocolCounts["Other"]},
    })

    return pie
}

func generateAnomalyChart() *charts.Bar {
    bar := charts.NewBar()
    bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
        Title: "AI Detected Anomalies",
    }))

    var anomalies []opts.BarData
    var times []string

    // Simulate adding detected anomalies for illustration
    for i := 0; i < 10; i++ {
        times = append(times, time.Now().Add(-time.Duration(i)*time.Minute).Format("15:04"))
        anomalies = append(anomalies, opts.BarData{Value: i})
    }

    bar.SetXAxis(times).
        AddSeries("Anomalies", anomalies)

    return bar
}

// RecordPacketCount logs the packet count per minute and updates protocol counts
func RecordPacketCount(packetProtocol string) {
    currentTime := time.Now().Format("15:04")
    packetCountMutex.Lock()
    packetCountPerMinute[currentTime]++
    protocolCounts[packetProtocol]++
    packetCountMutex.Unlock()
}

// SendAlert sends an alert message to all connected WebSocket clients
func SendAlert(alertMessage string) {
    alertsMutex.Lock()
    defer alertsMutex.Unlock()

    for _, conn := range alertsConn {
        if err := conn.WriteMessage(websocket.TextMessage, []byte(alertMessage)); err != nil {
            log.Println("Failed to send alert:", err)
            conn.Close()
        }
    }
}
