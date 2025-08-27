package postgres

import (
	"MockOrderService/internal/domain/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

// SaveOrder saves an order to the database
func (r *OrderRepository) SaveOrder(ctx context.Context, order *model.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
INSERT INTO orders (
  order_uid, track_number, entry, locale, internal_signature,
  customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
ON CONFLICT (order_uid) DO NOTHING
`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
ON CONFLICT (order_uid) DO NOTHING
`, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
		order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO payments (order_uid, transaction_id, request_id, currency, provider,
                      amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
ON CONFLICT (order_uid) DO NOTHING
`, order.OrderUID, order.Payment.TransactionID, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return err
	}

	// add batching for high load
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, `
INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
ON CONFLICT (order_uid, rid) DO NOTHING
`, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size,
			item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// GetOrderByOrderUID returns an order by orderUID from the database
func (r *OrderRepository) GetOrderByOrderUID(ctx context.Context, orderUID string) (*model.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var order model.Order
	err = tx.QueryRow(ctx,
		`SELECT order_uid, track_number, entry, locale, internal_signature,
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid = $1`, orderUID).
		Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.InternalSignature, &order.CustomerID,
			&order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		return nil, fmt.Errorf("orders query failed: %w", err)
	}

	var delivery model.Delivery
	err = tx.QueryRow(ctx,
		`SELECT order_uid, name, phone, zip, city, address, region, email 
			FROM deliveries WHERE order_uid = $1`, orderUID).
		Scan(&delivery.OrderUID, &delivery.Name, &delivery.Phone, &delivery.Zip,
			&delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
	if err != nil {
		return nil, fmt.Errorf("deliveries query failed: %w", err)
	}
	order.Delivery = &delivery

	rows, err := r.pool.Query(ctx,
		`SELECT id, order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, created_at
			 FROM items
			 WHERE order_uid = $1
	`, orderUID)
	if err != nil {
		return nil, fmt.Errorf("items failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		err = rows.Scan(&item.ID, &item.OrderUID, &item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status, &item.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("item scan failed: %w", err)
		}
		order.Items = append(order.Items, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("items iteration query failed: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetRecentOrders returns a slice of recent orders from the database
func (r *OrderRepository) GetRecentOrders(ctx context.Context, limit int) ([]*model.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var orders []*model.Order

	rows, err := r.pool.Query(ctx,
		`SELECT order_uid, track_number, entry, locale, internal_signature,
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("orders failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.InternalSignature, &order.CustomerID,
			&order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
		if err != nil {
			return nil, fmt.Errorf("order scan failed: %w", err)
		}
		orders = append(orders, &order)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("order iteration query failed: %w", err)
	}

	for _, order := range orders {
		var delivery model.Delivery
		err = tx.QueryRow(ctx,
			`SELECT order_uid, name, phone, zip, city, address, region, email 
			FROM deliveries WHERE order_uid = $1`, order.OrderUID).
			Scan(&delivery.OrderUID, &delivery.Name, &delivery.Phone, &delivery.Zip,
				&delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
		if err != nil {
			return nil, fmt.Errorf("deliveries query failed: %w", err)
		}
		order.Delivery = &delivery

		rows, err := r.pool.Query(ctx,
			`SELECT id, order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, created_at
			 FROM items
			 WHERE order_uid = $1
	`, order.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("items failed: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var item model.Item
			err = rows.Scan(&item.ID, &item.OrderUID, &item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status, &item.CreatedAt)
			if err != nil {
				return nil, fmt.Errorf("item scan failed: %w", err)
			}
			order.Items = append(order.Items, &item)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("items iteration query failed: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
