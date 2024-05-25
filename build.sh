#!/usr/bin/env bash

# Build the project
echo "Building for Windows x64"
GOOS=windows GOARCH=amd64 go build -o build/svdb-windows-x64.exe

echo "Building for Linux x64"
GOOS=linux GOARCH=amd64 go build -o build/svdb-linux-x64

echo "Building for MacOS x64"
GOOS=darwin GOARCH=amd64 go build -o build/svdb-macos-x64

echo "Building for Windows ARM64"
GOOS=windows GOARCH=arm64 go build -o build/svdb-windows-arm64.exe

echo "Building for Linux ARM64"
GOOS=linux GOARCH=arm64 go build -o build/svdb-linux-arm64

echo "Building for MacOS ARM"
GOOS=darwin GOARCH=arm64 go build -o build/svdb-macos-arm64

echo "Done"