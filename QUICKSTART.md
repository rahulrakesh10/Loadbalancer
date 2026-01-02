# Quick Start Guide

## 🚀 Getting Started in 3 Steps

### Step 1: Install Go
Make sure you have Go 1.21+ installed:
```bash
go version
```

### Step 2: Start Backend Servers

**Option A: Using the script (easiest)**
```bash
chmod +x scripts/start-backends.sh
./scripts/start-backends.sh
```

**Option B: Manual (3 separate terminals)**
```bash
# Terminal 1
go run cmd/backend/main.go -port=9001

# Terminal 2
go run cmd/backend/main.go -port=9002

# Terminal 3
go run cmd/backend/main.go -port=9003
```

### Step 3: Start Load Balancer

In a new terminal:
```bash
go run main.go
```

## ✅ Test It Works

```bash
# Make 3 requests - you should see different servers
curl http://localhost:8080
curl http://localhost:8080
curl http://localhost:8080

# Check metrics
curl http://localhost:8080/metrics
```

## 🐳 Or Use Docker (One Command)

```bash
docker-compose up --build
```

That's it! The load balancer is now running and distributing traffic.

## 📊 Expected Output

When you curl the load balancer multiple times, you should see:
```
Hello from server 9001
Hello from server 9002
Hello from server 9003
Hello from server 9001  (cycles back)
```

This confirms round-robin load balancing is working!

