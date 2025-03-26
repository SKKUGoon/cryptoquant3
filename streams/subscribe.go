package streams

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"cryptoquant.com/m/internal"
	"github.com/gorilla/websocket"
)

func SubscribeKline(symbol, interval string, handleKline func(internal.KlineData)) error {
	// Binance Futures WebSocket endpoint
	url := fmt.Sprintf("wss://fstream.binance.com/stream?streams=%s@kline_%s",
		strings.ToLower(symbol), interval)

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("websocket connection failed: %v", err)
	}
	defer conn.Close()

	log.Printf("Connected to Binance Futures stream for %s", symbol)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("error reading message: %v", err)
		}

		var streamData internal.Stream[internal.KlineData]
		if err := json.Unmarshal(message, &streamData); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Process the kline data using the provided handler
		handleKline(streamData.Data)
	}
}
