package data

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type BinanceFutureMarketData struct {
	RateLimit   int
	CurrentRate int
}

func NewBinanceFutureMarketData() *BinanceFutureMarketData {
	return &BinanceFutureMarketData{
		RateLimit:   1000,
		CurrentRate: 0,
	}
}

func (bm *BinanceFutureMarketData) UpdateRateLimit(rateLimit int) {
	bm.RateLimit = rateLimit
}

func (bm *BinanceFutureMarketData) UpdateCurrentRate(currentRate int) {
	bm.CurrentRate = currentRate
}

func (bm *BinanceFutureMarketData) GetKlineData(symbol string, interval string, limit int) {
	const weight = 5
	const urlBase = "https://fapi.binance.com/fapi/v1/klines"

	if bm.RateLimit-bm.CurrentRate < weight {
		// Wait until the start of next minute
		now := time.Now()
		nextMinute := now.Add(time.Minute).Truncate(time.Minute)
		time.Sleep(nextMinute.Sub(now))
		bm.CurrentRate = 0 // Reset rate limit counter for new minute
	}

	url := fmt.Sprintf("%s?symbol=%s&interval=%s&limit=%d", urlBase, symbol, interval, limit)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(body)
}
