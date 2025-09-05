// cmd/file-transfer/main.go
package main

import (
	"context"
	"log"
	
	"github.com/grandcat/zeroconf"
)

const serviceName = "_filetransfer._tcp"
const serviceport = 42424

func main() {
    server, err := zeroconf.Register(
		"GoFileTransferWindows",
		serviceName,
		"local.",
		serviceport,
	    nil,
		nil,
	)
	if err != nil {
		log.Fatalln("Failed to register service:", err)
	}
	defer server.Shutdown()

	log.Printf("Service %s registered on port %d", serviceName, serviceport)

	log.Printf("browsing for services...")
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)

	go func() {
		err = resolver.Browse(context.Background(), serviceName, "local.", entries)
		if err != nil {
			log.Fatalln("Failed to browse:", err)
		}
	}()

	log.Printf("Browsing for services...")

	go func() {
		for entry := range entries {
			if entry.Instance == "GoFileTransferWindows" {
				continue
			}
		    log.Printf("Found service: %s at %s:%d", entry.Instance, entry.AddrIPv4, entry.Port)
	    }
	}()
	
    select {}

}
