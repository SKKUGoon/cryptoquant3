package internal_test

import (
	"fmt"
	"testing"

	"cryptoquant.com/m/internal"
)

func TestNewFutureExchange(t *testing.T) {
	exchangeInfo, err := internal.NewFutureExchange()
	if err != nil {
		t.Fatalf("Failed to create exchange info: %v", err)
	}

	symbolInfo := exchangeInfo.GetSymbolInfo("BTCUSDT")
	if symbolInfo == nil {
		t.Fatalf("Failed to get symbol info: %v", err)
	}
}

func TestGetFutureSymbolFilter(t *testing.T) {
	exchangeInfo, err := internal.NewFutureExchange()
	if err != nil {
		t.Fatalf("Failed to create exchange info: %v", err)
	}

	symbolInfo := exchangeInfo.GetSymbolInfo("BTCUSDT")
	if symbolInfo == nil {
		t.Fatalf("Failed to get symbol info: %v", err)
	}

	symbolFilter := symbolInfo.GetSymbolFilter("LOT_SIZE")
	fmt.Println(symbolFilter)
}

func TestGetFutureSymbolPricePrecision(t *testing.T) {
	exchangeInfo, err := internal.NewFutureExchange()
	if err != nil {
		t.Fatalf("Failed to create exchange info: %v", err)
	}

	symbolInfo := exchangeInfo.GetSymbolInfo("BTCUSDT")
	if symbolInfo == nil {
		t.Fatalf("Failed to get symbol info: %v", err)
	}

	fmt.Println(symbolInfo.GetSymbolPricePrecision())
}
