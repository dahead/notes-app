#!/bin/sh
# For Windows
GOOS=windows GOARCH=amd64 go build -o notes-app.exe
mv notes-app.exe builds/windows/

# For Linux
GOOS=linux GOARCH=amd64 go build -o notes-app
