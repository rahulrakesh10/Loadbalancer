package balancer

import (
	"load-balancer/server"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// LoadBalancer handles HTTP request routing and load balancing
type LoadBalancer struct {
	rr          *RoundRobin
	healthCheck *HealthChecker
	proxy       *httputil.ReverseProxy
	mu          sync.RWMutex
}

// NewLoadBalancer creates a new load balancer instance
func NewLoadBalancer(backends []*server.BackendServer) *LoadBalancer {
	rr := NewRoundRobin(backends)
	
	// Start health checker (every 10 seconds, 5 second timeout)
	healthCheck := NewHealthChecker(backends, 10*time.Second, 5*time.Second)
	go healthCheck.Start()

	return &LoadBalancer{
		rr:          rr,
		healthCheck: healthCheck,
	}
}

// ServeHTTP implements http.Handler interface
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend, err := lb.rr.GetNextServer()
	if err != nil {
		http.Error(w, "No backend servers available", http.StatusServiceUnavailable)
		return
	}

	// Parse backend URL
	targetURL, err := url.Parse(backend.URL)
	if err != nil {
		log.Printf("Error parsing backend URL %s: %v", backend.URL, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create reverse proxy for this request
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Modify request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-For", r.RemoteAddr)
		req.Header.Set("X-Forwarded-Host", r.Host)
		req.Header.Set("X-Real-IP", r.RemoteAddr)
	}

	// Track connection
	backend.IncrementConn()
	defer backend.DecrementConn()

	// Log the request
	log.Printf("Routing request to %s", backend.URL)

	// Serve the request
	proxy.ServeHTTP(w, r)
}

// GetAliveBackends returns all alive backend servers
func (lb *LoadBalancer) GetAliveBackends() []*server.BackendServer {
	return lb.rr.GetAliveBackends()
}

