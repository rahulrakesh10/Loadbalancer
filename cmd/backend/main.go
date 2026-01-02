package main

import (
	"flag"
	"load-balancer/server"
	"log"
)

func main() {
	port := flag.Int("port", 9001, "Backend server port")
	flag.Parse()

	if *port < 1 || *port > 65535 {
		log.Fatalf("Invalid port number: %d", *port)
	}

	server.StartBackendServer(*port)
}


