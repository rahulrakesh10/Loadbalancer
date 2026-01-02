package balancer

import (
	"load-balancer/server"
	"log"
	"net/http"
	"time"
)

// HealthChecker performs periodic health checks on backend servers
type HealthChecker struct {
	backends []*server.BackendServer
	interval time.Duration
	timeout  time.Duration
	stopCh   chan struct{}
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(backends []*server.BackendServer, interval, timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		backends: backends,
		interval: interval,
		timeout:  timeout,
		stopCh:   make(chan struct{}),
	}
}

// Start begins the health checking process
func (hc *HealthChecker) Start() {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	log.Println("Health checker started")

	// Perform initial health check
	hc.checkAll()

	for {
		select {
		case <-ticker.C:
			hc.checkAll()
		case <-hc.stopCh:
			log.Println("Health checker stopped")
			return
		}
	}
}

// Stop stops the health checker
func (hc *HealthChecker) Stop() {
	close(hc.stopCh)
}

// checkAll performs health check on all backends
func (hc *HealthChecker) checkAll() {
	for _, backend := range hc.backends {
		go hc.checkHealth(backend)
	}
}

// checkHealth performs a health check on a single backend
func (hc *HealthChecker) checkHealth(backend *server.BackendServer) {
	client := &http.Client{
		Timeout: hc.timeout,
	}

	healthURL := backend.URL + "/health"
	resp, err := client.Get(healthURL)

	if err != nil {
		log.Printf("Health check failed for %s: %v", backend.URL, err)
		backend.SetAlive(false)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if !backend.IsAlive() {
			log.Printf("Backend %s is now healthy", backend.URL)
		}
		backend.SetAlive(true)
	} else {
		log.Printf("Health check failed for %s: status code %d", backend.URL, resp.StatusCode)
		backend.SetAlive(false)
	}
}


