package models

import (
	"time"
)

type Payment struct {
	ID                int64     `json:"id" db:"id"`
	UserID            int64     `json:"user_id" db:"user_id"`
	Amount            int64     `json:"amount" db:"amount"`
	Currency          string    `json:"currency" db:"currency"`
	ProductID         string    `json:"product_id" db:"product_id"`
	ProviderChargeID  string    `json:"provider_charge_id" db:"provider_charge_id"`
	TelegramChargeID  string    `json:"telegram_charge_id" db:"telegram_charge_id"`
	Status            string    `json:"status" db:"status"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"` // цена в копейках
	Currency    string `json:"currency"`
}

type PreCheckoutQuery struct {
	ID             string `json:"id"`
	TotalAmount    int64  `json:"total_amount"`
	Currency       string `json:"currency"`
	InvoicePayload string `json:"invoice_payload"`
}
