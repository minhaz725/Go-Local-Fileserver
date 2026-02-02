package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ip := ipnet.IP.To4()
			// Skip link-local addresses (169.254.x.x)
			if ip[0] == 169 && ip[1] == 254 {
				continue
			}
			return ipnet.IP.String()
		}
	}
	return "localhost (not connected to network)"
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

	http.HandleFunc("/ip", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "http://%s:%s", getLocalIP(), port)
	})

	http.Handle("/", noCacheMiddleware(http.FileServer(http.Dir(shareDir))))

	fmt.Printf("File server running on port %s\n", port)
	fmt.Printf("Local:   http://localhost:%s\n", port)
	fmt.Printf("Network: http://%s:%s\n", getLocalIP(), port)
	fmt.Printf("Current IP endpoint: /ip\n")
	fmt.Printf("Serving: %s\n", shareDir)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}