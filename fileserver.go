package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func getLocalIPs() []string {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ip := ipnet.IP.To4()
			// Skip link-local (169.254.x.x)
			if ip[0] == 169 && ip[1] == 254 {
				continue
			}
			// Skip VirtualBox (192.168.56.x)
			if ip[0] == 192 && ip[1] == 168 && ip[2] == 56 {
				continue
			}
			// Skip common VPN ranges (10.x.x.x)
			if ip[0] == 10 {
				continue
			}
			ips = append(ips, ipnet.IP.String())
		}
	}
	return ips
}

func noCacheMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		h.ServeHTTP(w, r)
	})
}

func main() {
	shareDir := "./shared"
	port := "8080"

	os.MkdirAll(shareDir, 0755)

	http.Handle("/", noCacheMiddleware(http.FileServer(http.Dir(shareDir))))

	ips := getLocalIPs()

	fmt.Println("========================================")
	fmt.Println("        Simple File Server")
	fmt.Println("========================================")
	fmt.Printf("Local:   http://localhost:%s\n", port)
	if len(ips) > 0 {
		fmt.Printf("Network: http://%s:%s\n", ips[0], port)
	} else {
		fmt.Println("Network: Not connected to local network")
	}
	fmt.Printf("Folder:  %s\n", shareDir)
	fmt.Println("========================================")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}