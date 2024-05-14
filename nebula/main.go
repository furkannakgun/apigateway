package main

import (
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/config"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/server"
)

func main() {
	// Load config
	config.LoadEnvVars()
	config.Load()
	// Start configuration watcher
	config.WatchConfig(server.ReloadRouter)

	// Initialize server
	server.Init()
	// Run
	server.Run()
}
