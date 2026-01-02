package balancer

import (
	"errors"
	"load-balancer/server"
	"sync"
)

// RoundRobin implements round-robin load balancing algorithm
type RoundRobin struct {
	backends []*server.BackendServer
	current  int
	mu       sync.Mutex
}

// NewRoundRobin creates a new round-robin load balancer
func NewRoundRobin(backends []*server.BackendServer) *RoundRobin {
	return &RoundRobin{
		backends: backends,
		current:  0,
	}
}

// GetNextServer returns the next available backend server using round-robin
func (rr *RoundRobin) GetNextServer() (*server.BackendServer, error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if len(rr.backends) == 0 {
		return nil, errors.New("no backends available")
	}

	// Try to find an alive server, starting from current position
	attempts := 0
	for attempts < len(rr.backends) {
		backend := rr.backends[rr.current]
		rr.current = (rr.current + 1) % len(rr.backends)

		if backend.IsAlive() {
			return backend, nil
		}

		attempts++
	}

	// If no alive server found, return the first one anyway (failover)
	if len(rr.backends) > 0 {
		return rr.backends[0], nil
	}

	return nil, errors.New("no backends available")
}

// GetAliveBackends returns all alive backends
func (rr *RoundRobin) GetAliveBackends() []*server.BackendServer {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	alive := make([]*server.BackendServer, 0)
	for _, backend := range rr.backends {
		if backend.IsAlive() {
			alive = append(alive, backend)
		}
	}
	return alive
}

// GetAllBackends returns all backends
func (rr *RoundRobin) GetAllBackends() []*server.BackendServer {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	return rr.backends
}


