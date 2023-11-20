package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	OrderID     uint64     `json:"order_id"`
	CustomerID  uuid.UUID  `json:"customer_id"`
	LineItems   []LineItem `json:"line_items"`
	CreateAt    *time.Time `json:"create_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type LineItem struct {
	ItemID   uuid.UUID `json:"item_id"`
	Quantity uint      `json:"quantity"`
	Price    uint      `json:"price"`
}
type FindResult struct {
	Orders []Order
	Cursor uint64
}

type ListPage struct {
	Size   uint64
	Offset uint64
}
