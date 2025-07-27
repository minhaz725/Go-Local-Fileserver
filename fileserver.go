package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    shareDir := "./shared"
    
    //create the directory if it doesn't exist
    if _, err := os.Stat(shareDir); os.IsNotExist(err) {
        os.Mkdir(shareDir, 0755)
        fmt.Printf("Created directory: %s\n", shareDir)
    }
    
    //file server handler
    fs := http.FileServer(http.Dir(shareDir))
    http.Handle("/", fs)
    
    port := "8080"
    fmt.Printf("File server running at:\n")
    fmt.Printf("Local: http://localhost:%s\n", port)
    fmt.Printf("Network: http://192.168.0.172:%s\n", port)
    fmt.Printf("Serving files from: %s\n", shareDir)
    
    log.Fatal(http.ListenAndServe(":"+port, nil))
}