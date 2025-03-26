package internal

// Keep all core structs here.

// PriceData represents real-time price information for a symbol
type PriceData struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Size   string `json:"size"`
	Time   string `json:"time"`
}

// TradeSignal represents a trading signal with direction
type TradeSignal struct {
	Pair      string `json:"pair"`
	Direction string `json:"direction"` // "long" or "short"
	Price     string `json:"price"`
	Time      string `json:"time"`
}

// Order represents a pair of long and short orders
type Order struct {
	LongOrder  SingleOrder `json:"long"`
	ShortOrder SingleOrder `json:"short"`
	Status     string      `json:"status"` // "open", "filled", "cancelled"
}

// SingleOrder represents an individual order with all necessary details
type SingleOrder struct {
	Symbol    string `json:"symbol"`
	Side      string `json:"side"` // "buy" or "sell"
	Size      string `json:"size"`
	Price     string `json:"price"`
	Time      string `json:"time"`
	OrderType string `json:"orderType"` // "market" or "limit"
	Status    string `json:"status"`    // "open", "filled", "cancelled"
}
