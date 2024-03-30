package model

import "time"

type CreditCard struct {
	ID        string    `db:"id"`
	Number    string    `db:"number"`
	Owner     string    `db:"owner"`
	ExpiresAt time.Time `db:"expires_at"`
	CVVCode   string    `db:"cvv_code"`
	PinCode   string    `db:"pin_code"`
	Metadata  string    `db:"metadata"`
}
