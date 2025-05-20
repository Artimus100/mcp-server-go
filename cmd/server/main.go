package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Artimus100/mcp-server-go/internal/config"

	"github.com/Artimus100/mcp-server-go/internal/handler"
	"github.com/Artimus100/mcp-server-go/internal/state"
	"github.com/Artimus100/mcp-server-go/internal/utils"
)

func main() {
	// Parse command line flags
	port := flag.Int("port", config.DefaultPort, "Port to listen on")
	flag.Parse()

	// Initialize logger
	logger := utils.NewLogger("server")
	logger.Info("Starting MCP server...")

	// Create context store
	contextStore := state.NewContextStore()

	// Create and start the server
	server := handler.NewServer(*port, contextStore, logger)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			logger.Error("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	logger.Info("MCP server listening on port %d", *port)

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-sigChan
	logger.Info("Received signal %v, shutting down...", sig)

	// Shutdown server
	if err := server.Shutdown(); err != nil {
		logger.Error("Error during shutdown: %v", err)
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}
