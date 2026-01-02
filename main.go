package main

import (
	"flag"
	"fmt"
	"load-balancer/balancer"
	"load-balancer/config"
	"load-balancer/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	// Parse command line flags
	configFile := flag.String("config", "config/servers.json", "Path to configuration file")
	lbPort := flag.Int("port", 8080, "Load balancer port")
	flag.Parse()

	// Load configuration
	var cfg *config.Config
	var err error

	if *configFile != "" {
		cfg, err = config.LoadConfig(*configFile)
		if err != nil {
			log.Printf("Failed to load config file %s: %v. Using default configuration.", *configFile, err)
			cfg = config.DefaultConfig()
		}
	} else {
		cfg = config.DefaultConfig()
	}

	// Override port if specified via flag
	if *lbPort != 8080 {
		cfg.Port = *lbPort
	}

	// Create backend servers from config
	backends := cfg.GetBackendServers()

	// Create load balancer
	lb := balancer.NewLoadBalancer(backends)

	// Setup HTTP server
	mux := http.NewServeMux()
	mux.Handle("/", lb)

	// Add metrics endpoint (optional)
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		aliveBackends := lb.GetAliveBackends()
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Load Balancer Metrics\n"))
		w.Write([]byte("====================\n\n"))
		w.Write([]byte("Alive Backends:\n"))
		for _, backend := range aliveBackends {
			w.Write([]byte("  - " + backend.URL + "\n"))
		}
	})

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: mux,
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Load balancer starting on port %d", cfg.Port)
		log.Printf("Backend servers configured: %d", len(backends))
		for _, backend := range backends {
			log.Printf("  - %s", backend.URL)
		}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start load balancer: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("\nShutting down load balancer...")
	// Note: In a production system, you'd want to properly stop the health checker
	// and wait for active connections to finish
}

