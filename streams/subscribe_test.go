package streams_test

import (
	"testing"

	"cryptoquant.com/m/streams"
)

func TestSubscribeKline(t *testing.T) {
	t.Log("Starting Stream BTCUSDT 1m")
	streams.SubscribeKline("BTCUSDT", "1m", func(kline streams.KlineData) {
		t.Log(kline)
	})
}
