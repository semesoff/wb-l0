package db

import (
	"database/sql"
	"errors"
	"fmt"
	"wb-l0/internal/models/order"
)

func (db *Database) AddOrder(order order.Order) error {
	DB := db.db
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	// add data into orders table
	_, err = tx.Exec(`
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderUid, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerId, order.DeliveryService, order.ShardKey, order.SmId, order.DateCreated, order.OofShard,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert into order: %v", err)
	}

	// add data into delivery table
	_, err = tx.Exec(`
		INSERT INTO delivery (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderUid, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert into delivery: %v", err)
	}

	// add data into payment table
	_, err = tx.Exec(`
		INSERT INTO payment (
			order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderUid, order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert into payment: %v", err)
	}

	for _, item := range order.Items {
		// add data into item table
		_, err = tx.Exec(`
			INSERT INTO item (
				order_uid, chrt_id, track_number, price, rid, name, sale, size,
				total_price, nm_id, brand, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT (order_uid) DO NOTHING`,
			order.OrderUid, item.ChrtId, item.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert into items: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (db *Database) GetOrder(orderUid string) (*order.Order, error) {
	DB := db.db
	var ord order.Order

	// get data from orders table
	err := DB.QueryRow(`
		SELECT
			order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders
		WHERE order_uid = $1
	`, orderUid).Scan(
		&ord.OrderUid, &ord.TrackNumber, &ord.Entry, &ord.Locale, &ord.InternalSignature,
		&ord.CustomerId, &ord.DeliveryService, &ord.ShardKey, &ord.SmId, &ord.DateCreated, &ord.OofShard,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("order not found: %s", orderUid)
		}
		return nil, fmt.Errorf("failed to query order: %v", err)
	}

	// get data from delivery table
	err = DB.QueryRow(`
		SELECT name, phone, zip, city, address, region, email
		FROM delivery
		WHERE order_uid = $1
	`, orderUid).Scan(
		&ord.Delivery.Name, &ord.Delivery.Phone, &ord.Delivery.Zip, &ord.Delivery.City,
		&ord.Delivery.Address, &ord.Delivery.Region, &ord.Delivery.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query delivery: %v", err)
	}

	// get data from payment table
	err = DB.QueryRow(`
		SELECT transaction, request_id, currency, provider, amount, payment_dt, bank,
			   delivery_cost, goods_total, custom_fee
		FROM payment
		WHERE order_uid = $1
	`, orderUid).Scan(
		&ord.Payment.Transaction, &ord.Payment.RequestId, &ord.Payment.Currency, &ord.Payment.Provider,
		&ord.Payment.Amount, &ord.Payment.PaymentDt, &ord.Payment.Bank,
		&ord.Payment.DeliveryCost, &ord.Payment.GoodsTotal, &ord.Payment.CustomFee,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query payment: %v", err)
	}

	// get data from item table
	rows, err := DB.Query(`
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price,
		       nm_id, brand, status
		FROM item
		WHERE order_uid = $1
	`, orderUid)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item order.Item
		err := rows.Scan(
			&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %v", err)
		}
		ord.Items = append(ord.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over items: %v", err)
	}

	return &ord, nil
}

// TODO: implement this method
func (db *Database) RestoreCache() {

}
