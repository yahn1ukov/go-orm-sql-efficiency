package models

type OrderProduct struct {
	ID        int
	OrderID   int
	ProductID int
	Product   Product
	Quantity  int
	Price     float64
}
