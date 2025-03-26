package streams_test

import (
	"testing"

	"cryptoquant.com/m/internal"
	"cryptoquant.com/m/streams"
)

func TestSubscribeKline(t *testing.T) {
	t.Log("Starting Stream BTCUSDT 1m")
	streams.SubscribeKline("BTCUSDT", "1m", func(kline internal.KlineData) {
		t.Log(kline)
	})
}
