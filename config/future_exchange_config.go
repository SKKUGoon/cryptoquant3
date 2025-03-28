package config

import (
	"sort"
	"strings"

	"cryptoquant.com/m/internal"
)

// FutureTradeConfig implements the Exchange interface
type FutureTradeConfig struct {
	ExchangeInfo *internal.FutureExchange

	// List of symbols to exclude from trading
	ExcludeTrades map[string]bool
	// Trading Parameters
	MaximumLeverage int
	LossLimit       float64
	ProfitLimit     float64

	// List of symbols to subscribe to
	QuotingAsset string

	isTest bool
}

func NewFutureTradeConfig() (*FutureTradeConfig, error) {
	if exchangeInfo, err := internal.NewFutureExchange(); err != nil {
		return nil, err
	} else {
		return &FutureTradeConfig{
			ExchangeInfo: exchangeInfo,
		}, nil
	}
}

func (e *FutureTradeConfig) SetTestMode(value bool) {
	e.isTest = value
}

func (e *FutureTradeConfig) UpdateExchangeInfo() {
	if exchangeInfo, err := internal.NewFutureExchange(); err != nil {
		panic(err)
	} else {
		e.ExchangeInfo = exchangeInfo
	}
}

func (e *FutureTradeConfig) UpdateMaximumLeverage(value int) {
	e.MaximumLeverage = value
}

func (e *FutureTradeConfig) UpdateLossLimit(value float64) {
	e.LossLimit = value
}

func (e *FutureTradeConfig) UpdateProfitLimit(value float64) {
	e.ProfitLimit = value
}

func (e *FutureTradeConfig) UpdateQuotingAsset(value string) {
	e.QuotingAsset = value
}

func (e *FutureTradeConfig) GetSymbolQuotePrecision(symbol string) int {
	return e.ExchangeInfo.GetSymbolInfo(symbol).GetSymbolPricePrecision()
}

func (e *FutureTradeConfig) GetAvailableSymbols() []string {
	symbols := e.ExchangeInfo.GetAvailableSymbols(e.isTest)
	quotingSymbols := make([]string, 0)
	for _, symbol := range symbols {
		if strings.HasSuffix(symbol, e.QuotingAsset) {
			quotingSymbols = append(quotingSymbols, symbol)
		}
	}
	return quotingSymbols
}

// CreatePair implements Exchange interface
func (e *FutureTradeConfig) CreatePair(test bool) []string {
	symbols := e.ExchangeInfo.GetAvailableSymbols(test)
	sort.Strings(symbols) // ensures deterministic order
	pairs := make([]string, len(symbols)*(len(symbols)-1)/2)
	k := 0
	for i := range symbols {
		if _, ok := e.ExcludeTrades[symbols[i]]; ok {
			continue
		}

		for j := i + 1; j < len(symbols); j++ {
			if _, ok := e.ExcludeTrades[symbols[j]]; ok {
				continue
			}

			// Add the pair to the list (e.g. BTCUSDT-ETHUSDT)
			pairs[k] = symbols[i] + "-" + symbols[j]
			k++
		}
	}
	return pairs
}

func (e *FutureTradeConfig) GetSymbolPricePrecision(symbol string) int {
	return e.ExchangeInfo.GetSymbolInfo(symbol).GetSymbolPricePrecision()
}
