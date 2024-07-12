#!/bin/bash

# Create the releases directory if it doesn't exist
mkdir -p releases

# Build for different platforms
GOOS=darwin GOARCH=amd64 go build -o releases/fabrun-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o releases/fabrun-darwin-arm64
GOOS=linux GOARCH=amd64 go build -o releases/fabrun-linux-amd64
GOOS=windows GOARCH=amd64 go build -o releases/fabrun-windows-amd64.exe

echo "Build completed and binaries are placed in the 'releases' directory."