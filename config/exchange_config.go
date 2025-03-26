package config

import "cryptoquant.com/m/internal"

// Exchange holds all the configuration for the exchange
type Exchange struct {
	ExchangeInfo *internal.Exchange
	SymbolInfo   *internal.SymbolInfo
	SymbolFilter *internal.SymbolFilter
}

// NewExchange creates a new exchange
func NewExchange(exchangeInfo *internal.Exchange) *Exchange {
	return &Exchange{
		ExchangeInfo: exchangeInfo,
	}
}

// GetSymbolInfo returns the symbol information for the given symbol
func (e *Exchange) GetSymbolInfo(symbol string) *internal.SymbolInfo {
	return e.ExchangeInfo.GetSymbolInfo(symbol)
}

// GetSymbolFilter returns the filter of specified type for the given symbol
func (e *Exchange) GetSymbolFilter(symbol string, filterType string) *internal.SymbolFilter {
	return e.SymbolInfo.GetSymbolFilter(filterType)
}
