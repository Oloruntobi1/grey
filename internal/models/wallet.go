package models

import (
	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID      string          `json:"id"`
	UserID  string          `json:"user_id"`
	Balance decimal.Decimal `json:"balance"`
}
