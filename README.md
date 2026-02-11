# Custom HTTP Load Balancer in Go

A high-performance HTTP load balancer built in Go that distributes incoming client requests across multiple backend servers using configurable load-balancing algorithms. Features health checks, concurrency-safe request routing, and fault tolerance.

## 🚀 Features

- **Round-Robin Load Balancing**: Distributes requests evenly across backend servers
- **Health Checks**: Periodic health monitoring with automatic failover
- **Concurrency-Safe**: Thread-safe server selection using mutexes
- **Reverse Proxy**: Full HTTP reverse proxy functionality
- **Fault Tolerance**: Automatically skips unhealthy servers
- **Configurable**: JSON-based configuration for backend servers
- **Metrics Endpoint**: Monitor load balancer status
- **Docker Support**: Containerized deployment ready

## 📋 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (optional, for containerized deployment)

## 🏗️ Architecture

```
Client → Load Balancer (:8080) → Backend Servers (:9001, :9002, :9003)
```

The load balancer:
1. Receives HTTP requests on port 8080
2. Selects a healthy backend server using round-robin algorithm
3. Forwards the request to the selected backend
4. Returns the response to the client

## 📁 Project Structure

```
load-balancer/
├── main.go                 # Load balancer entry point
├── cmd/
│   └── backend/
│       └── main.go         # Backend server entry point
├── server/
│   └── backend.go          # Backend server implementation
├── balancer/
│   ├── round_robin.go      # Round-robin algorithm
│   ├── health_check.go     # Health checking logic
│   └── proxy.go            # Reverse proxy implementation
├── config/
│   ├── config.go           # Configuration loader
│   └── servers.json        # Server configuration
├── scripts/
│   └── start-backends.sh   # Helper script to start backends
├── Dockerfile              # Docker image definition
├── docker-compose.yml      # Docker Compose configuration
└── README.md
```

## 🛠️ Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd load-balancer
```

2. Install dependencies:
```bash
go mod download
```

## 🚦 Usage

### Option 1: Manual Setup (Recommended for Development)

#### Step 1: Start Backend Servers

You can start backend servers manually in separate terminals:

**Terminal 1:**
```bash
go run cmd/backend/main.go -port=9001
```

**Terminal 2:**
```bash
go run cmd/backend/main.go -port=9002
```

**Terminal 3:**
```bash
go run cmd/backend/main.go -port=9003
```

Or use the helper script:
```bash
chmod +x scripts/start-backends.sh
./scripts/start-backends.sh
```

#### Step 2: Start Load Balancer

In a new terminal:
```bash
go run main.go
```

Or with custom configuration:
```bash
go run main.go -config=config/servers.json -port=8080
```

## 🔍 How It Works

### Round-Robin Algorithm

The load balancer uses a round-robin algorithm to distribute requests:
- Request 1 → Server 1
- Request 2 → Server 2
- Request 3 → Server 3
- Request 4 → Server 1 (cycles back)

