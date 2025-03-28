package streams_test

import (
	"testing"
	"time"

	"cryptoquant.com/m/internal"
	"cryptoquant.com/m/streams"
)

func TestSubscribeKline(t *testing.T) {
	t.Log("Starting Stream BTCUSDT 1m")
	ch := make(chan internal.KlineDataStream)
	done := make(chan struct{})
	go streams.SubscribeKline("BTCUSDT", "1m", ch, done)

	// Send done signal after 10 seconds
	go func() {
		time.Sleep(10 * time.Second)
		done <- struct{}{}
	}()

	for kline := range ch {
		t.Log(kline)
	}
}

func TestSubscribeKlineMulti(t *testing.T) {
	t.Log("Starting Stream BTCUSDT 1m")
	chMap := make(map[string]chan internal.KlineDataStream)
	chMap["BTCUSDT"] = make(chan internal.KlineDataStream)
	chMap["ETHUSDT"] = make(chan internal.KlineDataStream)
	chMap["XRPUSDT"] = make(chan internal.KlineDataStream)
	chMap["DOGEUSDT"] = make(chan internal.KlineDataStream)
	chMap["LINKUSDT"] = make(chan internal.KlineDataStream)

	symbols := []string{"BTCUSDT", "ETHUSDT", "XRPUSDT", "DOGEUSDT", "LINKUSDT"}
	done := make(chan struct{})
	go streams.SubscribeKlineMulti(symbols, "1m", chMap, done)

	// Send done signal after 10 seconds
	go func() {
		time.Sleep(10 * time.Second)
		done <- struct{}{}
	}()

	for {
		select {
		case <-done:
			return
		case kline := <-chMap["BTCUSDT"]:
			t.Log(kline)
		case kline := <-chMap["ETHUSDT"]:
			t.Log(kline)
		case kline := <-chMap["XRPUSDT"]:
			t.Log(kline)
		case kline := <-chMap["DOGEUSDT"]:
			t.Log(kline)
		case kline := <-chMap["LINKUSDT"]:
			t.Log(kline)
		}
	}
}
