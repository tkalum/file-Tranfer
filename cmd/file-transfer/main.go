//cmd/file-transfer/main.go
package main


import (
	"file-transfer/internal/discovery"
	"file-transfer/internal/transfer"
	"log"
	"os"
	"flag"
	"context"
	"fmt"
	"net"

)

func sendMode(args []string) {
	filename := args[0]

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", discovery.ServicePort))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "UnknownDevice"
	}

	announce, err := discovery.AnnounceService(hostname)
	if err != nil {
		log.Fatal(err)
	}
	defer announce.Shutdown()
	log.Printf("Service announced")

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ip := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	log.Printf("Connection established with %s", ip)
    
	transfer.Dialer(filename, ip)
	
	log.Println("File sent successfully")
}

func getMode() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transferDone := make(chan error, 1)
	go func() {
		err := transfer.Listener(transfer.TransferPort)
		transferDone <- err
	}()


	devices, err := discovery.BrowseServices(ctx)
	if err != nil {
		log.Fatalf("❌ Browse error: %v", err)
	}
	var firstdevice *discovery.Device

	for device := range devices {
		firstdevice = device
		log.Printf("🔍 Discovered: %s at %v:%d", device.Instance, device.IP, device.Port)
		cancel()
		break
	}
	if firstdevice == nil {
		log.Println("No devices found")
		return
	}

	log.Printf("Connecting to %s at %v:%d", firstdevice.Instance, firstdevice.IP, firstdevice.Port)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", firstdevice.IP.String(), firstdevice.Port))
	if err != nil {
		log.Fatalf("❌ Connection error: %v", err)
	}
	defer conn.Close()
	log.Printf("Connected to %s", firstdevice.Instance)

	err = transfer.Listener(discovery.ServicePort)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	//listnmode := flag.Bool("l", false, "Run as listener")
	send := flag.String("s", "", "File to send")
	flag.Parse()

	if *send != "" {
		// Pass the value from the port flag to the Listener function.
		sendMode([]string{*send})

	} else {
		
		getMode()
    } 
} 