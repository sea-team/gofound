#!/usr/bin/env bash

version=`cat version.txt`

echo "Building..."
echo "删除旧的文件"
rm -rf ./dist

echo "构建darwin-amd4"
GOOS=darwin GOARCH=amd64 go build -o ./dist/macOS/macos_amd64_$version

echo "构建darwin-arm64(Apple Silicon)"
GOOS=darwin GOARCH=arm64 go build -o ./dist/macOS/macos_apple_silicon_$version

echo "构建linux-amd64"
GOOS=linux GOARCH=amd64 go build -o ./dist/linux/linx_amd64_$version

echo "构建linux-arm64"
GOOS=linux GOARCH=arm64 go build -o ./dist/linux/linx_arm64_$version

echo "构建windows-amd4"
GOOS=windows GOARCH=amd64 go build -o ./dist/windows/windows_amd64_$version.exe

echo "构建windows-arm64"
GOOS=windows GOARCH=arm64 go build -o ./dist/windows/windows_arm64_$version.exe