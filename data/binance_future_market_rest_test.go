package data_test

import (
	"testing"

	"cryptoquant.com/m/data"
)

func TestGetKlineData(t *testing.T) {
	bm := &data.BinanceFutureMarketData{
		RateLimit: 1000,
	}
	t.Log("Usage", bm.CurrentRate)
	bm.GetKlineData("BTCUSDT", "1m", 100)
	t.Log("Usage", bm.CurrentRate)
}
