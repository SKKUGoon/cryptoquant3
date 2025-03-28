package config_test

import (
	"fmt"
	"testing"

	"cryptoquant.com/m/config"
)

func TestFutureTradeConfig_CreatePair(t *testing.T) {
	tradeConfig, err := config.NewFutureTradeConfig()
	if err != nil {
		t.Fatal(err)
	}

	pair := tradeConfig.CreatePair(true)

	fmt.Println(pair)
}
