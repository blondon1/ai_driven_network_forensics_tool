package ui

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strings"
)

var (
	username = "admin"
	password = "password"
)

func StartServer() {
	http.HandleFunc("/", basicAuth(trafficChartHandler))
	http.HandleFunc("/charts", basicAuth(trafficChartHandler))
	http.HandleFunc("/alerts", basicAuth(alertsHandler))
	http.HandleFunc("/historical", basicAuth(historicalDataHandler))

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

// basicAuth is a middleware function that provides basic HTTP authentication
func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !validateCredentials(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// validateCredentials checks if the provided username and password are correct
func validateCredentials(user, pass string) bool {
	return user == username && pass == password
}

func trafficChartHandler(w http.ResponseWriter, r *http.Request) {
	page := components.NewPage()
	page.AddCharts(generateTrafficLineChart())
	page.Render(w)
}

func historicalDataHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		query = "SELECT * FROM packets ORDER BY timestamp DESC LIMIT 100"
	}

	db, err := sql.Open("sqlite3", "data/packets.db")
	if err != nil {
		log.Println("Failed to open database:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Println("Failed to execute query:", err)
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result strings.Builder
	for rows.Next() {
		var id int
		var timestamp, sourceIP, destinationIP, protocol string
		var size int
		err := rows.Scan(&id, &timestamp, &sourceIP, &destinationIP, &protocol, &size)
		if err != nil {
			log.Println("Failed to scan row:", err)
			http.Error(w, "Row scan error", http.StatusInternalServerError)
			return
		}
		result.WriteString(fmt.Sprintf("%s %s -> %s %s %d bytes\n",
			timestamp, sourceIP, destinationIP, protocol, size))
	}
	w.Write([]byte(result.String()))
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
