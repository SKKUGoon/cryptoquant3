package main

import (
	"fmt"

	"cryptoquant.com/m/internal"
)

func main() {
	pd := internal.PriceData{
		Symbol: "BTCUSDT",
		Price:  "10000",
		Size:   "1",
		Time:   "1715817600",
	}
	fmt.Println(pd)
}
