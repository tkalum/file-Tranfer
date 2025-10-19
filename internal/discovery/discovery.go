// internal/discovery/discovery.go
package discovery

import (
	"context"
	"log"
	"net"
	
	"github.com/grandcat/zeroconf"
)

const ServiceName = "_filetransfer._tcp"
const ServicePort = 24242

type Device struct {
	Instance string
	IP net.IP
	Port int
}


func AnnounceService(instance string) (*zeroconf.Server, error) {
	server, err := zeroconf.Register(
		instance,
		ServiceName,
		"local.",
		ServicePort,
		nil,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
	return server, nil	
}

func BrowseServices(ctx context.Context) (<-chan *Device, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}
	entries := make(chan *zeroconf.ServiceEntry)
	devices := make(chan *Device)

	go func() {
		defer close(devices)
		for entry := range entries {
			log.Printf("Discovered service: %s at %s:%d", entry.Instance, entry.AddrIPv4, entry.Port)
			devices <- &Device{
				Instance: entry.Instance,
				IP: entry.AddrIPv4[0],
				Port: entry.Port,
			}
		}
	}()

	err = resolver.Browse(ctx, ServiceName, "local.", entries)
	if err != nil {
		return nil, err
	}

	return devices, nil
}