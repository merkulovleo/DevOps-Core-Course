package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var startTime = time.Now()

type ServiceInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type SystemInfo struct {
	Hostname     string `json:"hostname"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
	NumCPU       int    `json:"num_cpu"`
	GoVersion    string `json:"go_version"`
}

type RuntimeInfo struct {
	UptimeSeconds float64 `json:"uptime_seconds"`
	CurrentTime   string  `json:"current_time"`
	Timezone      string  `json:"timezone"`
}

type Response struct {
	Service ServiceInfo `json:"service"`
	System  SystemInfo  `json:"system"`
	Runtime RuntimeInfo `json:"runtime"`
}

type HealthResponse struct {
	Status        string  `json:"status"`
	Timestamp     string  `json:"timestamp"`
	UptimeSeconds float64 `json:"uptime_seconds"`
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	log.Printf("GET / requested from %s", r.RemoteAddr)

	response := Response{
		Service: ServiceInfo{
			Name:        "DevOps Info Service (Go)",
			Version:     "1.0.0",
			Description: "A web service providing system and runtime information",
		},
		System: SystemInfo{
			Hostname:     getHostname(),
			Platform:     runtime.GOOS,
			Architecture: runtime.GOARCH,
			NumCPU:       runtime.NumCPU(),
			GoVersion:    runtime.Version(),
		},
		Runtime: RuntimeInfo{
			UptimeSeconds: time.Since(startTime).Seconds(),
			CurrentTime:   time.Now().UTC().Format(time.RFC3339),
			Timezone:      "UTC",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET /health requested from %s", r.RemoteAddr)

	response := HealthResponse{
		Status:        "healthy",
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		UptimeSeconds: time.Since(startTime).Seconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/health", healthHandler)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Starting DevOps Info Service (Go) on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
