package main

import (
	"log"
	"time"

	"cryptoquant.com/m/engine"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	engine := engine.NewEngine()
	engine.StartStreamCh()
	engine.StartAssets()
	engine.StartPairs()
	engine.StartStream()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Printf("[engine status] %s", engine.GetStatus())
	}
}
