package config

import "cryptoquant.com/m/internal"

// SpotExchange holds all the configuration for the exchange
type SpotExchange struct {
	ExchangeInfo *internal.SpotExchange
	SymbolInfo   *internal.SpotSymbolInfo
	SymbolFilter *internal.SpotSymbolFilter
}

// NewSpotExchange creates a new exchange
func NewSpotExchange(exchangeInfo *internal.SpotExchange) *SpotExchange {
	return &SpotExchange{
		ExchangeInfo: exchangeInfo,
	}
}

// GetSymbolInfo returns the symbol information for the given symbol
func (e *SpotExchange) GetSymbolInfo(symbol string) *internal.SpotSymbolInfo {
	return e.ExchangeInfo.GetSymbolInfo(symbol)
}

// GetSymbolFilter returns the filter of specified type for the given symbol
func (e *SpotExchange) GetSymbolFilter(symbol string, filterType string) *internal.SpotSymbolFilter {
	return e.SymbolInfo.GetSymbolFilter(filterType)
}
