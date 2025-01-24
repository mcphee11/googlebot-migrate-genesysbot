#!/bin/bash

go build -o linux/mcpheeware-migrate && env GOOS=windows GOARCH=amd64 go build -o windows/mcpheeware-migrate.exe && env GOOS=darwin GOARCH=amd64 go build -o macos/mcpheeware-migrate