package config

// TradeConfig defines the common interface for both spot and futures exchanges
type TradeConfig interface {
	GetSymbolInfo(symbol string) interface{} // returns either *internal.SymbolInfo or *internal.FutureSymbolInfo
	CreatePair() []string
}
