package ui

import (
    "fmt"
    "github.com/go-echarts/go-echarts/v2/charts"
    "github.com/go-echarts/go-echarts/v2/components"
    "github.com/go-echarts/go-echarts/v2/opts"
    "github.com/gorilla/websocket"
    "html/template"
    "log"
    "net/http"
    "sync"
    "time"
)

var (
    alertsMutex       sync.Mutex
    alertsConn        []*websocket.Conn
    upgrader          = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        CheckOrigin:     func(r *http.Request) bool { return true },
    }
    packetCountMutex     sync.Mutex
    packetCountPerMinute = make(map[string]int)
    alertChannel         = make(chan string)
)

func StartServer() {
    http.HandleFunc("/", dashboardHandler)
    http.HandleFunc("/charts", trafficChartHandler)
    http.HandleFunc("/alerts", alertsHandler)
    http.HandleFunc("/historical", historicalDataHandler)

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
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
    })

    go func() {
        for alert := range alertChannel {
            BroadcastAlert(alert)
        }
    }()

    log.Println("Starting web server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    page := components.NewPage()
    page.AddCharts(generateTrafficLineChart(), generateTrafficBarChart())
    page.Render(w)
}

func trafficChartHandler(w http.ResponseWriter, r *http.Request) {
    page := components.NewPage()
    page.AddCharts(generateTrafficLineChart())
    page.Render(w)
}

func alertsHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Real-time Alerts will be displayed here")
}

func historicalDataHandler(w http.ResponseWriter, r *http.Request) {
    tpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Historical Data</title>
    </head>
    <body>
        <h1>Historical Data</h1>
        <table border="1">
            <tr>
                <th>Timestamp</th>
                <th>Source IP</th>
                <th>Destination IP</th>
                <th>Protocol</th>
                <th>Size</th>
            </tr>
            {{range .}}
            <tr>
                <td>{{.Timestamp}}</td>
                <td>{{.SourceIP}}</td>
                <td>{{.DestinationIP}}</td>
                <td>{{.Protocol}}</td>
                <td>{{.Size}}</td>
            </tr>
            {{end}}
        </table>
    </body>
    </html>`

    // Example data, replace this with actual data retrieval
    data := []struct {
        Timestamp     string
        SourceIP      string
        DestinationIP string
        Protocol      string
        Size          int
    }{
        {"2024-08-12 10:00:00", "192.168.1.1", "192.168.1.2", "TCP", 1500},
        {"2024-08-12 10:01:00", "192.168.1.2", "192.168.1.1", "UDP", 512},
    }

    t, err := template.New("table").Parse(tpl)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    t.Execute(w, data)
}

func generateTrafficLineChart() *charts.Line {
    line := charts.NewLine()
    line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
        Title: "Network Traffic",
    }))

    packetCountMutex.Lock()
    defer packetCountMutex.Unlock()

    // Prepare data
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

func generateTrafficBarChart() *charts.Bar {
    bar := charts.NewBar()
    bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
        Title: "Network Traffic Distribution",
    }))

    packetCountMutex.Lock()
    defer packetCountMutex.Unlock()

    // Prepare data
    var times []string
    var counts []opts.BarData
    for timeStr, count := range packetCountPerMinute {
        times = append(times, timeStr)
        counts = append(counts, opts.BarData{Value: count})
    }

    bar.SetXAxis(times).
        AddSeries("Packets per Minute", counts)

    return bar
}

// RecordPacketCount logs the packet count per minute
func RecordPacketCount() {
    currentTime := time.Now().Format("15:04")
    packetCountMutex.Lock()
    packetCountPerMinute[currentTime]++
    packetCountMutex.Unlock()
}

// BroadcastAlert sends an alert message to all connected WebSocket clients
func BroadcastAlert(message string) {
    alertsMutex.Lock()
    defer alertsMutex.Unlock()

    for _, conn := range alertsConn {
        err := conn.WriteMessage(websocket.TextMessage, []byte(message))
        if err != nil {
            log.Println("Failed to send alert:", err)
            conn.Close()
        }
    }
}

// SendAlert sends an alert to the alertChannel
func SendAlert(alert string) {
    alertChannel <- alert
}
