#!/bin/bash

# Script to start multiple backend servers

echo "Starting backend servers..."

# Start backend servers in background
go run cmd/backend/main.go -port=9001 &
PID1=$!
echo "Backend server 1 started on port 9001 (PID: $PID1)"

go run cmd/backend/main.go -port=9002 &
PID2=$!
echo "Backend server 2 started on port 9002 (PID: $PID2)"

go run cmd/backend/main.go -port=9003 &
PID3=$!
echo "Backend server 3 started on port 9003 (PID: $PID3)"

echo ""
echo "All backend servers started!"
echo "Press Ctrl+C to stop all servers"

# Wait for interrupt
trap "kill $PID1 $PID2 $PID3; exit" INT TERM
wait


