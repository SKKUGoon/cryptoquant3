package main

import (
	"fmt"
	"log"

	"cryptoquant.com/m/internal"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	pd := internal.PriceData{
		Symbol: "BTCUSDT",
		Price:  "10000",
		Size:   "1",
		Time:   "1715817600",
	}
	fmt.Println(pd)
}
