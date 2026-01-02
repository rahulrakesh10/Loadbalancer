package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// BackendServer represents a backend server in the pool
type BackendServer struct {
	URL          string
	Alive        bool
	ActiveConns  int
	LastChecked  time.Time
	mu           sync.RWMutex
}

// NewBackendServer creates a new backend server instance
func NewBackendServer(url string) *BackendServer {
	return &BackendServer{
		URL:         url,
		Alive:       true,
		ActiveConns: 0,
		LastChecked: time.Now(),
	}
}

// SetAlive sets the alive status of the server
func (b *BackendServer) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Alive = alive
	b.LastChecked = time.Now()
}

// IsAlive returns the current alive status
func (b *BackendServer) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Alive
}

// IncrementConn increments the active connection count
func (b *BackendServer) IncrementConn() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ActiveConns++
}

// DecrementConn decrements the active connection count
func (b *BackendServer) DecrementConn() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.ActiveConns > 0 {
		b.ActiveConns--
	}
}

// GetActiveConns returns the current active connection count
func (b *BackendServer) GetActiveConns() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.ActiveConns
}

// StartBackendServer starts a simple HTTP server on the specified port
func StartBackendServer(port int) {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Hello from server %d\n", port)
	})
	
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})
	
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Backend server starting on port %d", port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Failed to start backend server on port %d: %v", port, err)
	}
}

