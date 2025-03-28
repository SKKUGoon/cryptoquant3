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
	// Check the weight usage and how it's filling up
	for i := 0; i < 10; i++ {
		klineData, err := bm.GetKlineData("BTCUSDT", "1m", 100)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(klineData)
		t.Log("Usage", bm.CurrentRate, "/", bm.RateLimit, "used")
	}
}
