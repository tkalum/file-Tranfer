// internal/discovery/discovery.go
package discovery

import (
	"context"
	"log"
	"net"
	"strings"
	
	"github.com/grandcat/zeroconf"
)

const ServiceName = "_filetransfer._tcp"
const ServicePort = 24242

type Device struct {
	Instance string
	IP net.IP
	Port int
	Filename string
	Filesize string
}

func findLocalIP(prefix string) net.IP {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	
	for _, iface := range ifaces {
		// Skip interfaces that are down or are loopback (127.0.0.1)
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue 
		}
		
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		
		for _, addr := range addrs {
			var ip net.IP
			// Extract IP from address structure
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			
			// Only consider IPv4 addresses and check for the desired Wi-Fi prefix
			if ip != nil && ip.To4() != nil && strings.HasPrefix(ip.String(), prefix) {
				return ip
			}
		}
	}
	return nil
}


func AnnounceService(instance string, Filename string, Filesize string) (*zeroconf.Server, error) {
	localIP := findLocalIP("192.168.") 
    if localIP == nil {
        log.Println("❌ Error: Could not find a local IP on the 192.168.8.x network. Check WiFi connection.")
        return nil, nil
    }
	log.Printf("✅ Local IP found: %s\n", localIP.String())

	matadata := []string{
		"Filename=" + Filename,
		"Filesize=" + Filesize,
	}
	server, err := zeroconf.RegisterProxy(
		instance,
		ServiceName,
		"local.",
		ServicePort,
		localIP.String(),
		[]string{localIP.String()},
		matadata,
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
			
			mataData := make(map[string]string)
			for _, txt := range entry.Text {
				parts := strings.SplitN(txt, "=", 2)
				if len(parts) == 2 {
					mataData[parts[0]] = parts[1]
				}
			}

			filename := mataData["Filename"]
			filesize := mataData["Filesize"]

			devices <- &Device{
				Instance: entry.Instance,
				IP: entry.AddrIPv4[0],
				Port: entry.Port,
				Filename: filename,
				Filesize: filesize,
			}
		}
	}()

	err = resolver.Browse(ctx, ServiceName, "local.", entries)
	if err != nil {
		return nil, err
	}

	return devices, nil
}