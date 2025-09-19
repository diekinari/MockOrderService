package model

import "time"

type Order struct {
	OrderUID          string     `json:"order_uid" validate:"required"`
	TrackNumber       string     `json:"track_number,omitempty" validate:"required"`
	Entry             string     `json:"entry,omitempty"`
	Locale            string     `json:"locale,omitempty" validate:"omitempty,max=10"`
	InternalSignature string     `json:"internal_signature,omitempty"`
	CustomerID        string     `json:"customer_id,omitempty"`
	DeliveryService   string     `json:"delivery_service,omitempty"`
	Shardkey          string     `json:"shardkey,omitempty"`
	SmID              *int32     `json:"sm_id,omitempty" validate:"omitempty,gte=0"`
	DateCreated       *time.Time `json:"date_created,omitempty"` // custom check: not in future (validate manually)
	OofShard          string     `json:"oof_shard,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`

	// delivery/payment required by business logic -> required tag ensures non-nil
	Delivery *Delivery `json:"delivery,omitempty" validate:"required"`
	Payment  *Payment  `json:"payment,omitempty" validate:"required"`

	// items must contain at least one element; dive => apply element-level tags
	Items []*Item `json:"items,omitempty" validate:"min=1,dive,required"`
}

type Delivery struct {
	ID        int32      `json:"id"`
	OrderUID  string     `json:"order_uid"`
	Name      string     `json:"name,omitempty" validate:"required"`
	Phone     string     `json:"phone,omitempty" validate:"required,phone"` // phone — custom validator registered (e.g. +7XXXXXXXXXX)
	Zip       string     `json:"zip,omitempty" validate:"required,ziprus"`  // ziprus — custom validator (6 digits)
	City      string     `json:"city,omitempty" validate:"required"`
	Address   string     `json:"address,omitempty" validate:"required"`
	Region    string     `json:"region,omitempty"`
	Email     string     `json:"email,omitempty" validate:"omitempty,email"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type Payment struct {
	ID            int32      `json:"id"`
	OrderUID      string     `json:"order_uid"`
	TransactionID string     `json:"transaction_id,omitempty"`
	RequestID     string     `json:"request_id,omitempty"`
	Currency      string     `json:"currency,omitempty"`
	Provider      string     `json:"provider,omitempty"`
	Amount        *int64     `json:"amount,omitempty" validate:"required,gte=0"`     // pointer + required ensures non-nil
	PaymentDt     *int64     `json:"payment_dt,omitempty" validate:"omitempty,gt=0"` // custom: also check not-future if needed
	Bank          string     `json:"bank,omitempty"`
	DeliveryCost  *int64     `json:"delivery_cost,omitempty" validate:"omitempty,gte=0"`
	GoodsTotal    *int64     `json:"goods_total,omitempty" validate:"omitempty,gte=0"`
	CustomFee     *int64     `json:"custom_fee,omitempty" validate:"omitempty,gte=0"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
}

type Item struct {
	ID          int32      `json:"id"`
	OrderUID    string     `json:"order_uid"`
	ChrtID      *int64     `json:"chrt_id,omitempty" validate:"omitempty"`
	TrackNumber string     `json:"track_number,omitempty"`
	Price       *int64     `json:"price,omitempty" validate:"required,gte=0"`
	Rid         string     `json:"rid,omitempty" validate:"required"`
	Name        string     `json:"name,omitempty"`
	Sale        *int32     `json:"sale,omitempty" validate:"omitempty,gte=0,lte=100"`
	Size        string     `json:"size,omitempty"`
	TotalPrice  *int64     `json:"total_price,omitempty" validate:"required,gte=0"`
	NmID        *int64     `json:"nm_id,omitempty" validate:"omitempty"`
	Brand       string     `json:"brand,omitempty"`
	Status      *int32     `json:"status,omitempty" validate:"omitempty,gte=0"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}
