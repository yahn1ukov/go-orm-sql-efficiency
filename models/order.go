package models

import "time"

type Order struct {
	ID         int
	CustomerID int
	Date       time.Time
	Total      float64
	Products   []OrderProduct
	CreatedAt  time.Time
}
