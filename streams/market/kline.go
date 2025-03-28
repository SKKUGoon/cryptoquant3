package streams

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"cryptoquant.com/m/internal"
	"github.com/gorilla/websocket"
)

func SubscribeKline(symbol, interval string, ch chan internal.KlineDataStream, done chan struct{}) error {
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

	var streamData internal.Stream[internal.KlineDataStream]
	for {
		select {
		case <-done:
			close(ch)
			return nil
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				return fmt.Errorf("error reading message: %v", err)
			}

			if err := json.Unmarshal(message, &streamData); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// Process the kline data using the provided handler
			ch <- streamData.Data
		}
	}
}

func SubscribeKlineMulti(symbols []string, interval string, chMap map[string]chan internal.KlineDataStream, done chan struct{}) error {
	urlBase := "wss://fstream.binance.com/stream?streams="
	streams := make([]string, len(symbols))
	for i := range symbols {
		symbol := symbols[i]
		streams[i] = fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval)
	}

	url := urlBase + strings.Join(streams, "/")
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("websocket connection failed: %v", err)
	}
	defer conn.Close()

	log.Printf("Connected to Binance Futures stream for %s", strings.Join(symbols, ", "))

	var streamData internal.Stream[internal.KlineDataStream]
	for {
		select {
		case <-done:
			// Close all channels
			for _, symbol := range symbols {
				close(chMap[symbol])
			}
			err := conn.WriteMessage(websocket.CloseMessage, []byte{})
			if err != nil {
				log.Printf("Error closing connection: %v", err)
			}
			println("Connection to Binance Futures stream for %s is closed", strings.Join(symbols, ", "))
			return nil
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				return fmt.Errorf("error reading message: %v", err)
			}

			if err := json.Unmarshal(message, &streamData); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			chMap[streamData.Data.Symbol] <- streamData.Data
		}
	}
}
