package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ipni/herald/herald"
	"github.com/ipni/go-libipni/metadata"
)

func main() {
	log.Println("Starting application...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Creating Metadata instance...")

	// Create a Metadata instance
	mc := metadata.Default.(metadata.MetadataContext)

	// For this example, let's assume you want to use the Bitswap transport protocol.
	// You might need to initialize Bitswap or other transport protocols with necessary data.
	bitswapProtocol := &metadata.Bitswap{}  // Initialize with appropriate values if needed
	log.Println("Bitswap protocol initialized.")

	md := mc.New(bitswapProtocol)
	log.Println("Metadata instance created.")

	// Initialize the Herald with desired options
	log.Println("Initializing Herald...")
	h, err := herald.New(
		herald.WithMetadata(md),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Herald: %v", err)
	}
	log.Println("Herald initialized successfully.")

	// Start Herald
	log.Println("Starting Herald...")
	if err := h.Start(ctx); err != nil {
		log.Fatalf("Failed to start Herald: %v", err)
	}
	log.Println("Herald started.")

	// Set up signal catching
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Signal catching set up. Waiting for shutdown signal...")

	// Wait for a signal to shutdown
	<-sigChan

	// Shutdown gracefully with a timeout
	log.Println("Shutdown signal received. Gracefully shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := h.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Failed to gracefully shutdown Herald: %v", err)
	}

	log.Println("Herald service stopped.")
}
