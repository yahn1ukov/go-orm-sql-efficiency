package models

import "time"

type Customer struct {
	ID        int
	Name      string
	Email     string
	Orders    []Order
	CreatedAt time.Time
}
