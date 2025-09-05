// cmd/file-transfer/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grandcat/zeroconf"
)

const serviceName = "_filetransfer._tcp"
const servicePort = 42424

func main() {
	// Get hostname for instance name
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "UnknownDevice"
	}

	// Step 1: Start advertising service
	server, err := zeroconf.Register(
		hostname,
		serviceName,
		"local.",
		servicePort,
		nil,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
	defer server.Shutdown()
	log.Printf("✅ Service registered: %s on port %d", hostname, servicePort)

	// Step 2: Start TCP server to listen for connections
	go startTCPServer()

	// Step 3: Discover other services
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalf("Failed to create resolver: %v", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	
	// Start browsing for services
	go func() {
		ctx := context.Background()
		err := resolver.Browse(ctx, serviceName, "local.", entries)
		if err != nil {
			log.Printf("Browse error: %v", err)
		}
	}()

	// Step 4: Process discovered services
	go func() {
		for entry := range entries {
			// Skip self-discovery
			if entry.Instance == hostname {
				continue
			}
			
			log.Printf("🔍 Discovered: %s at %v:%d", 
				entry.Instance, entry.AddrIPv4, entry.Port)
			
			// Try to connect to discovered device
			if len(entry.AddrIPv4) > 0 {
				go connectToDevice(entry.AddrIPv4[0].String(), entry.Port)
			}
		}
	}()

	// Step 5: Manual connection option
	go func() {
		time.Sleep(5 * time.Second) // Wait a bit before offering manual connect
		fmt.Printf("\n💡 Want to connect manually? Enter IP address (or press Enter to skip): ")
		var ip string
		fmt.Scanln(&ip)
		if ip != "" {
			connectToDevice(ip, servicePort)
		}
	}()

	// Step 6: Wait for exit signal
	log.Printf("🚀 Running! Press Ctrl+C to exit")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")
}

// Start TCP server to listen for incoming connections
func startTCPServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", servicePort))
	if err != nil {
		log.Fatalf("❌ Failed to start TCP server: %v", err)
	}
	defer listener.Close()
	
	log.Printf("👂 Listening for connections on port %d", servicePort)
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("❌ Connection accept error: %v", err)
			continue
		}
		
		log.Printf("✅ Incoming connection from %s", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

// Handle incoming connection
func handleConnection(conn net.Conn) {
	defer conn.Close()
	
	// Send welcome message
	conn.Write([]byte(fmt.Sprintf("Hello from %s!\n", getHostname())))
	
	// Here you would add your file transfer logic
	log.Printf("📁 File transfer session started with %s", conn.RemoteAddr())
	
	// Keep connection open for demo
	time.Sleep(30 * time.Minute) // Or implement actual file transfer
}

// Connect to another device manually
func connectToDevice(ip string, port int) {
	log.Printf("🚀 Attempting connection to %s:%d", ip, port)
	
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 5*time.Second)
	if err != nil {
		log.Printf("❌ Connection to %s:%d failed: %v", ip, port, err)
		return
	}
	defer conn.Close()
	
	log.Printf("✅ Connected to %s:%d", ip, port)
	
	// Read welcome message
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	log.Printf("📩 Message from remote: %s", string(buffer[:n]))
	
	// Here you would add your file transfer logic
}

// Get hostname helper
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "UnknownDevice"
	}
	return hostname
}