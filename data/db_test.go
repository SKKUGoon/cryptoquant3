package data_test

import (
	"fmt"
	"testing"

	"cryptoquant.com/m/data"
	"github.com/joho/godotenv"
)

func TestImportNimbusData(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	db, err := data.ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	prices1, prices2, err := db.ImportNimbusData()
	if err != nil {
		t.Fatalf("Failed to import Nimbus data: %v", err)
	}

	fmt.Println(prices1)
	fmt.Println(prices2)
}
