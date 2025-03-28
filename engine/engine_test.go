package engine_test

import (
	"testing"
	"time"

	"cryptoquant.com/m/engine"
	"github.com/joho/godotenv"
)

func TestEngine_StartStreamCh(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	engine := engine.NewEngine()
	engine.SetTestMode(true)
	engine.FutureConfig.UpdateQuotingAsset("USDT")

	engine.StartStreamCh()

	t.Log(engine.ChMap)
}

func TestEngine_StartPairAssets(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// done := make(chan struct{})

	engine := engine.NewEngine()
	engine.SetTestMode(true)
	engine.FutureConfig.UpdateQuotingAsset("USDT")

	engine.StartStreamCh()
	engine.StartAssets()

	t.Log(engine.PairAssets)
}

func TestEngine_StartPairs(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	engine := engine.NewEngine()
	engine.SetTestMode(true)
	engine.FutureConfig.UpdateQuotingAsset("USDT")

	engine.StartStreamCh()
	engine.StartAssets()
	engine.StartPairs()

	t.Log(engine.Pairs)
}

func TestEngine_StartStream(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	engine := engine.NewEngine()
	engine.SetTestMode(true)
	engine.FutureConfig.UpdateQuotingAsset("USDT")

	engine.StartStreamCh()
	engine.StartAssets()
	engine.StartPairs()
	engine.StartStream()

	time.Sleep(26 * time.Second)

	for _, pair := range engine.Pairs {
		t.Log(pair.PairSymbol, pair.HistoricSpread.Items)
	}
}
