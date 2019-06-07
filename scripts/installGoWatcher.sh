#!/bin/bash

# Get the package with:
go get github.com/canthefason/go-watcher

# Install the binary under go/bin folder:
go install github.com/canthefason/go-watcher/cmd/watcher

# Make sure that your go/bin folder is appended to PATH environment variable.

# Then run: watcher
# If you come across segmentation error, set up the GOPATH env var first!
# GOPATH=~/go watcher
