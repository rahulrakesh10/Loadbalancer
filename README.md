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

### Option 2: Docker Compose (Recommended for Production)

Start everything with Docker Compose:
```bash
docker-compose up --build
```

This will start:
- 3 backend servers (ports 9001, 9002, 9003)
- 1 load balancer (port 8080)

## 🧪 Testing

### Basic Test with curl

Make multiple requests to see load balancing in action:

```bash
curl http://localhost:8080
curl http://localhost:8080
curl http://localhost:8080
```

You should see responses rotating:
```
Hello from server 9001
Hello from server 9002
Hello from server 9003
```

### Check Metrics

View load balancer metrics:
```bash
curl http://localhost:8080/metrics
```

### Load Testing with Apache Bench

```bash
ab -n 1000 -c 10 http://localhost:8080/
```

### Load Testing with hey

```bash
hey -n 1000 -c 10 http://localhost:8080/
```

## ⚙️ Configuration

Edit `config/servers.json` to add or modify backend servers:

```json
{
  "backends": [
    {
      "url": "http://localhost:9001"
    },
    {
      "url": "http://localhost:9002"
    },
    {
      "url": "http://localhost:9003"
    }
  ],
  "port": 8080
}
```

## 🔍 How It Works

### Round-Robin Algorithm

The load balancer uses a round-robin algorithm to distribute requests:
- Request 1 → Server 1
- Request 2 → Server 2
- Request 3 → Server 3
- Request 4 → Server 1 (cycles back)

The algorithm is thread-safe using mutexes to handle concurrent requests.

### Health Checks

- Health checks run every 10 seconds
- Each backend is checked via `/health` endpoint
- Unhealthy servers are automatically skipped
- Failed servers are retried on the next health check cycle

### Concurrency

- Each HTTP request is handled in its own goroutine
- Shared state (server selection) is protected with mutexes
- Supports thousands of concurrent connections

## 📊 Key Concepts Demonstrated

1. **Concurrency**: Goroutines and channels for concurrent request handling
2. **Thread Safety**: Mutexes for protecting shared state
3. **Networking**: HTTP reverse proxy using `net/http/httputil`
4. **Fault Tolerance**: Health checks and automatic failover
5. **System Design**: Load balancing architecture and patterns



## 🔧 Advanced Features

### Custom Health Check Interval

Modify the health check interval in `balancer/proxy.go`:

```go
healthCheck := NewHealthChecker(backends, 10*time.Second, 5*time.Second)
//                                 ^interval  ^timeout
```

### Adding More Backend Servers

1. Add server configuration to `config/servers.json`
2. Start the new backend server:
```bash
go run cmd/backend/main.go -port=9004
```

### Implementing Other Algorithms

The codebase is structured to easily add new load balancing algorithms:
- Least Connections
- Weighted Round Robin
- Random Selection

## 🐛 Troubleshooting


