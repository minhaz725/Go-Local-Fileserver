# Simple Go File Server

Basic HTTP file server for sharing files on home network

## Usage
- Run: go run fileserver.go
- Add files to shared/ folder  
- Access from browser at displayed URL
- Logs saved to fileserver.log

## Build
Windows: go build -o fileserver.exe fileserver.go
Mac: set GOOS=darwin && set GOARCH=amd64 && go build -o fileserver-mac fileserver.go

## What it does
- Serves files from shared/ directory
- Shows directory listings
- Works cross-platform