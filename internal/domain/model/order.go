package model

import "time"

type Order struct {
	OrderUID          string     `json:"order_uid"`
	TrackNumber       string     `json:"track_number,omitempty"`
	Entry             string     `json:"entry,omitempty"`
	Locale            string     `json:"locale,omitempty"`
	InternalSignature string     `json:"internal_signature,omitempty"`
	CustomerID        string     `json:"customer_id,omitempty"`
	DeliveryService   string     `json:"delivery_service,omitempty"`
	Shardkey          string     `json:"shardkey,omitempty"`
	SmID              *int32     `json:"sm_id,omitempty"`
	DateCreated       *time.Time `json:"date_created,omitempty"`
	OofShard          string     `json:"oof_shard,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`

	Delivery *Delivery `json:"delivery,omitempty"`
	Payment  *Payment  `json:"payment,omitempty"`
	Items    []*Item   `json:"items,omitempty"`
}

type Delivery struct {
	ID        int32      `json:"id"`
	OrderUID  string     `json:"order_uid"`
	Name      string     `json:"name,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	Zip       string     `json:"zip,omitempty"`
	City      string     `json:"city,omitempty"`
	Address   string     `json:"address,omitempty"`
	Region    string     `json:"region,omitempty"`
	Email     string     `json:"email,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type Payment struct {
	ID            int32      `json:"id"`
	OrderUID      string     `json:"order_uid"`
	TransactionID string     `json:"transaction_id,omitempty"`
	RequestID     string     `json:"request_id,omitempty"`
	Currency      string     `json:"currency,omitempty"`
	Provider      string     `json:"provider,omitempty"`
	Amount        *int64     `json:"amount,omitempty"`
	PaymentDt     *int64     `json:"payment_dt,omitempty"`
	Bank          string     `json:"bank,omitempty"`
	DeliveryCost  *int64     `json:"delivery_cost,omitempty"`
	GoodsTotal    *int64     `json:"goods_total,omitempty"`
	CustomFee     *int64     `json:"custom_fee,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
}

type Item struct {
	ID          int32      `json:"id"`
	OrderUID    string     `json:"order_uid"`
	ChrtID      *int64     `json:"chrt_id,omitempty"`
	TrackNumber string     `json:"track_number,omitempty"`
	Price       *int64     `json:"price,omitempty"`
	Rid         string     `json:"rid,omitempty"`
	Name        string     `json:"name,omitempty"`
	Sale        *int32     `json:"sale,omitempty"`
	Size        string     `json:"size,omitempty"`
	TotalPrice  *int64     `json:"total_price,omitempty"`
	NmID        *int64     `json:"nm_id,omitempty"`
	Brand       string     `json:"brand,omitempty"`
	Status      *int32     `json:"status,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}
