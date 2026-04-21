package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=10.95.17.149 user=omniadm password=0mn14dm dbname=omniport port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	type Result struct {
		ID      uint64
		MenuText string
		MenuURL  string
		Table    string
	}

	var results []Result
	db.Raw("SELECT id, menu_text, menu_url, 'menus' as table FROM adm.posm_menus WHERE menu_url LIKE '%infra%'").Scan(&results)
	
	var accessResults []Result
	db.Raw("SELECT 0 as id, menu_text, menu_url, 'access' as table FROM adm.posm_access WHERE menu_url LIKE '%infra%'").Scan(&accessResults)
	
	results = append(results, accessResults...)

	fmt.Println("=== Infrastructure Results ===")
	for _, r := range results {
		fmt.Printf("[%v] [%v] %v -> %v\n", r.Table, r.ID, r.MenuText, r.MenuURL)
	}
}
