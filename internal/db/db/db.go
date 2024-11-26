package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"wb-l0/config"
	"wb-l0/internal/models/order"
)

type DatabaseProvider interface {
	AddOrder(order order.Order) error
	GetOrder(orderUid string) (*order.Order, error)
}

type Database struct {
	db *sql.DB
}

func NewDatabase(cfg *config.Database) *Database {
	db := &Database{}
	db.InitDB(cfg)
	return db
}

func (DB *Database) InitDB(cfg *config.Database) {
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
	DB.db = db
}

//func (DB *Database) GetDB() *Database {
//	if DB == nil || DB.db == nil {
//		log.Fatal("Database is not initialized. Call InitDB() first.")
//		return nil
//	}
//	return DB
//}
