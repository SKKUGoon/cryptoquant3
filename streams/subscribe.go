package streams

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type KlineData struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	Kline     struct {
		StartTime    int64  `json:"t"`
		CloseTime    int64  `json:"T"`
		Symbol       string `json:"s"`
		Interval     string `json:"i"`
		FirstTradeID int64  `json:"f"`
		LastTradeID  int64  `json:"L"`
		OpenPrice    string `json:"o"`
		ClosePrice   string `json:"c"`
		HighPrice    string `json:"h"`
		LowPrice     string `json:"l"`
		Volume       string `json:"v"`
		NumTrades    int    `json:"n"`
		IsClosed     bool   `json:"x"`
		QuoteVolume  string `json:"q"`
	} `json:"k"`
}

func SubscribeKline(symbol, interval string, handleKline func(KlineData)) error {
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

		var kline KlineData
		if err := json.Unmarshal(message, &kline); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Process the kline data using the provided handler
		handleKline(kline)
	}
}

// func Streams(symbol string, url string) {
// 	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
// 	if err != nil {
// 		log.Fatal("Failed to connect to WebSocket:", err)
// 	}
// 	defer conn.Close()

// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Fatal("Failed to read message:", err)
// 		}

// 		price :=
// 	}
// }
