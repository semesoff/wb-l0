package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"wb-l0/config"
)

type Database struct {
	db *sql.DB
}

var DB *Database

func InitDB() {
	cfg := config.GetConfig().Database
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name))
	if err != nil {
		log.Fatalf("Error openning database: %v", err)
		return
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
		return
	}
	DB = &Database{db: db}
}

func GetDB() *sql.DB {
	if DB == nil || DB.db == nil {
		log.Fatal("Database is not initialized. Call InitDB() first.")
		return nil
	}
	return DB.db
}

func CloseDB() {
	if DB != nil && DB.db != nil {
		if err := DB.db.Close(); err != nil {
			log.Fatalf("Error closing database: %v", err)
			return
		}
	} else {
		log.Fatalf("Database is not initialized.")
		return
	}
	log.Println("Database is closed.")
}
