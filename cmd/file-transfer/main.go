//cmd/file-transfer/main.go
package main

import (
	"fmt"
	"log"
	"net"

	"github.com/grandcat/zeroconf"

)

const serviceName = "_filetransfer._tcp"
const serviceport = 42424

func main() {
    server, err ;= zeroconf.Register(
		"GoFileTransfer",
		serviceName,
		"local.",
		serviceport,
	    nil,
		nil
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
	defer resolver.Close()

	
}
