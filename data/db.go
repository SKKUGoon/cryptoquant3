package data

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strconv"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func ConnectDB() (*Database, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		os.Getenv("USERNAME"),
		os.Getenv("PASSWORD"),
		os.Getenv("PG_NAME"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

// ImportNimbusData imports Nimbus data for testing purposes
func (d *Database) ImportNimbusData() ([]float64, []float64, error) {
	query := `SELECT timestamp, close, close_time, symbol, quote_volume 
		FROM nimbus.binance_klines 
		WHERE (symbol = 'APEUSDT' or symbol = 'BTTCUSDT') 
		AND timestamp BETWEEN $1 AND $2`

	timestamp1 := 1739718000000
	timestamp2 := 1740322799000

	rows, err := d.db.Query(query, timestamp1, timestamp2)
	if err != nil {
		return nil, nil, fmt.Errorf("error querying data: %v", err)
	}
	defer rows.Close()

	// Use a map to group data by symbol
	dataBySymbol := make(map[string][]struct {
		CloseTime int64
		Close     float64
	})

	for rows.Next() {
		var (
			timestamp int64
			close     string
			closeTime int64
			symbol    string
			quoteVol  float64
		)

		if err := rows.Scan(&timestamp, &close, &closeTime, &symbol, &quoteVol); err != nil {
			return nil, nil, fmt.Errorf("error scanning row: %v", err)
		}

		closePrice, err := strconv.ParseFloat(close, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing close price: %v", err)
		}

		dataBySymbol[symbol] = append(dataBySymbol[symbol], struct {
			CloseTime int64
			Close     float64
		}{
			CloseTime: closeTime,
			Close:     closePrice,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating rows: %v", err)
	}

	// Sort each symbol's data by CloseTime
	for symbol, data := range dataBySymbol {
		sort.Slice(data, func(i, j int) bool {
			return data[i].CloseTime < data[j].CloseTime
		})

		dataBySymbol[symbol] = data
	}

	// Extracting close prices
	var prices1, prices2 []float64
	for symbol, data := range dataBySymbol {
		for _, d := range data {
			if symbol == "APEUSDT" {
				prices1 = append(prices1, d.Close)
			} else if symbol == "BTTCUSDT" {
				prices2 = append(prices2, d.Close)
			}
		}
	}

	return prices1, prices2, nil
}

// ExportTradeLog exports trade log data
func (d *Database) ExportTradeLog() error {
	// Export trade log data
	return nil
}
