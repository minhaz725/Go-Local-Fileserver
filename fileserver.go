package main

import (
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "strings"
)

// Get the local IP address of the machine
func getLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "localhost"
    }

    for _, addr := range addrs {
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                // Skip virtual interfaces like VirtualBox, VMware, etc.
                if !strings.Contains(ipnet.IP.String(), "169.254") && // Skip link-local
                   !strings.Contains(ipnet.IP.String(), "192.168.56") { // Skip VirtualBox
                    return ipnet.IP.String()
                }
            }
        }
    }
    return "localhost"
}

func main() {
    shareDir := "./shared"
    
    // Create the directory if it doesn't exist
    if _, err := os.Stat(shareDir); os.IsNotExist(err) {
        os.Mkdir(shareDir, 0755)
        fmt.Printf("Created directory: %s\n", shareDir)
    }
    
    // File server handler
    fs := http.FileServer(http.Dir(shareDir))
    http.Handle("/", fs)
    
    port := "8080"
    localIP := getLocalIP()
    
    fmt.Printf("File server running at:\n")
    fmt.Printf("Local: http://localhost:%s\n", port)
    fmt.Printf("Network: http://%s:%s\n", localIP, port)
    fmt.Printf("Serving files from: %s\n", shareDir)
    
    log.Fatal(http.ListenAndServe(":"+port, nil))
}