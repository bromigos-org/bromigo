package main

import (
	"github.com/bromigos-org/bromigo/common/run"
)

func main() {
	run.StartHTTPServer() // Start the HTTP server for health checks
	run.Init()
}
