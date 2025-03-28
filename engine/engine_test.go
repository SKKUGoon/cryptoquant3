package engine_test

import (
	"fmt"
	"testing"
	"time"

	"cryptoquant.com/m/engine"
	"github.com/joho/godotenv"
)

func TestEngine(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	engine := engine.NewEngine()
	engine.SetTestMode(true)

	engine.StartChMap()

	t.Log(engine.ChMap)
}

func TestEngineStartPairCalculation(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	engine := engine.NewEngine()
	engine.SetTestMode(true)
	engine.FutureConfig.UpdateQuotingAsset("USDT")

	engine.StartChMap()
	engine.PreparePairCalculation()

	t.Log(engine)
}

func TestEngineStartStream(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	done := make(chan struct{})

	engine := engine.NewEngine()
	engine.SetTestMode(true)
	engine.FutureConfig.UpdateQuotingAsset("USDT")

	engine.StartChMap()
	engine.PreparePairCalculation()
	engine.ListenPair(done)
	engine.StartStream(done)

	time.Sleep(10 * time.Second)
	done <- struct{}{}

	for key, i := range engine.Pairs {
		fmt.Println(key, i.Spread)
	}
	fmt.Println(engine.Pairs)
}
