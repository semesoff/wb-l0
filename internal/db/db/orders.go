package db

import (
	_ "errors"
	"fmt"
	"log"
	"wb-l0/internal/cache"
	"wb-l0/internal/models/order"
	"wb-l0/internal/utils"
)

func (db *Database) AddOrder(order order.Order, dataBytes []byte) error {
	DB := db.db

	// add data into orders table
	_, err := DB.Exec(`INSERT INTO orders (order_uid, order_data) VALUES ($1, $2) ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderUid, dataBytes)
	if err != nil {
		return fmt.Errorf("failed to insert into order: %v", err)
	}
	return nil
}

func (db *Database) GetOrder(orderUid string) (*order.Order, error) {
	DB := db.db

	// get data from orders table
	var dataBytes []byte
	row := DB.QueryRow(`SELECT * FROM orders WHERE order_uid = $1`, orderUid)
	err := row.Scan(&orderUid, &dataBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %v", err)
	}

	// bytes to model
	orderData, err := utils.EncodeMessage(dataBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode message: %v", err)
	}
	return &orderData, nil
}

func (db *Database) RestoreCache(rp *cache.RedisProvider) {
	DB := db.db
	rows, err := DB.Query(`SELECT * FROM orders`)
	if err != nil {
		log.Fatalf("failed to get orders: %v", err)
		return
	}

	for rows.Next() {
		var orderUid string
		var dataBytes []byte
		err := rows.Scan(&orderUid, &dataBytes)
		if err != nil {
			log.Printf("failed to scan order: %v", err)
			return
		}
		if err := (*rp).SetCache(orderUid, dataBytes); err != nil {
			log.Printf("failed to add order to cache: %v", err)
			return
		}
	}
	log.Println("Cache is restored.")
}
