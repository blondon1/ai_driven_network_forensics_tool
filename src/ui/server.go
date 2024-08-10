package ui

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	packetCountPerMinute = make(map[string]int)
	packetCountMutex     sync.Mutex
)

func StartServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/charts", http.StatusFound)
	})

	http.HandleFunc("/charts", trafficChartHandler)
	http.HandleFunc("/alerts", alertsHandler)

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

	log.Println("Starting web server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func trafficChartHandler(w http.ResponseWriter, r *http.Request) {
	page := components.NewPage()
	page.AddCharts(generateTrafficLineChart())
	page.Render(w)
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
		AddSeries("Packets per Minute", counts).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	return line
}

// RecordPacketCount logs the packet count per minute
func RecordPacketCount() {
	currentTime := time.Now().Format("15:04")
	packetCountMutex.Lock()
	packetCountPerMinute[currentTime]++
	packetCountMutex.Unlock()
}
