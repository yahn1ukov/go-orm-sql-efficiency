package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/yahn1ukov/go-orm-sql-efficiency/models"
	"github.com/yahn1ukov/go-orm-sql-efficiency/utils"
)

type SQL struct {
	db *sql.DB
}

type CreateProduct struct {
	SQL
}

func (op *CreateProduct) Name() string {
	return "Create Product"
}

func (op *CreateProduct) Execute(iteration int) error {
	_, err := op.db.
		Exec(
			"INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4)",
			fmt.Sprintf("Product_%d", iteration),
			"Test product",
			99.99,
			100,
		)

	return err
}

type GetCustomerByID struct {
	ID int
	SQL
}

func (op *GetCustomerByID) Name() string {
	return "Get Customer by ID"
}

func (op *GetCustomerByID) Execute(int) error {
	var customer models.Customer
	return op.db.
		QueryRow(
			"SELECT id, name, email, created_at FROM customers WHERE id = $1",
			op.ID,
		).
		Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.CreatedAt,
		)
}

type UpdateProductByID struct {
	ID    int
	Price float64
	SQL
}

func (op *UpdateProductByID) Name() string {
	return "Update Product Price by ID"
}

func (op *UpdateProductByID) Execute(int) error {
	_, err := op.db.
		Exec(
			"UPDATE products SET price = $1 WHERE id = $2",
			op.Price,
			op.ID,
		)

	return err
}

type DeleteProductByName struct {
	SQL
}

func (op *DeleteProductByName) Name() string {
	return "Delete Product by Name"
}

func (op *DeleteProductByName) Execute(iteration int) error {
	_, err := op.db.Exec(
		"DELETE FROM products WHERE name = $1",
		fmt.Sprintf("Product_%d", iteration),
	)

	return err
}

type CreateOrderWithProductsByCustomerID struct {
	CustomerID int
	ProductID  int
	SQL
}

func (op *CreateOrderWithProductsByCustomerID) Name() string {
	return "Create Order with Products By Customer ID (Transaction)"
}

func (op *CreateOrderWithProductsByCustomerID) Execute(int) error {
	tx, err := op.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var orderID uint
	err = tx.
		QueryRow(
			"INSERT INTO orders (customer_id, total) VALUES ($1, $2) RETURNING id",
			op.CustomerID,
			199.98,
		).
		Scan(&orderID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO order_products (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)",
		orderID,
		op.ProductID,
		2,
		99.99,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

type GetCustomerStatsByID struct {
	CustomerID int
	SQL
}

func (op *GetCustomerStatsByID) Name() string {
	return "Get Customer Stats by ID (Aggregation)"
}

func (op *GetCustomerStatsByID) Execute(int) error {
	var result struct {
		TotalOrders int64
		TotalSpent  float64
	}

	return op.db.
		QueryRow(
			`SELECT 
            COUNT(*) AS total_orders, 
            COALESCE(SUM(total), 0) AS total_spent 
        FROM orders 
        WHERE customer_id = $1`,
			op.CustomerID,
		).
		Scan(
			&result.TotalOrders,
			&result.TotalSpent,
		)
}

type GetProductSalesByLimit struct {
	Limit int
	SQL
}

func (op *GetProductSalesByLimit) Name() string {
	return "Get Product Sales by Limit (Complex Join)"
}

func (op *GetProductSalesByLimit) Execute(int) error {
	rows, err := op.db.
		Query(
			`SELECT 
            product_id, 
            products.name AS product_name, 
            SUM(quantity) AS total_sales, 
            SUM(order_products.price * quantity) AS revenue
        FROM order_products
        JOIN products ON products.id = order_products.product_id
        GROUP BY product_id, products.name
        ORDER BY revenue DESC
        LIMIT $1`,
			op.Limit,
		)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var r struct {
			ProductID   uint
			ProductName string
			TotalSales  int64
			Revenue     float64
		}
		if err := rows.Scan(
			&r.ProductID,
			&r.ProductName,
			&r.TotalSales,
			&r.Revenue,
		); err != nil {
			return err
		}
	}

	return rows.Err()
}

type GetOrderFullDetailsByID struct {
	OrderID int
	SQL
}

func (op *GetOrderFullDetailsByID) Name() string {
	return "Get Order Full Details by ID (Nested Preload)"
}

func (op *GetOrderFullDetailsByID) Execute(int) error {
	var order models.Order
	err := op.db.
		QueryRow(
			"SELECT id, customer_id, date, total, created_at FROM orders WHERE id = $1",
			op.OrderID,
		).
		Scan(
			&order.ID,
			&order.CustomerID,
			&order.Date,
			&order.Total,
			&order.CreatedAt,
		)
	if err != nil {
		return err
	}

	rows, err := op.db.
		Query(
			`SELECT 
            op.id, op.order_id, op.product_id, op.quantity, op.price,
            p.name, p.description, p.price, p.stock, p.created_at
        FROM order_products op
        JOIN products p ON p.id = op.product_id
        WHERE op.order_id = $1`,
			op.OrderID,
		)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var op models.OrderProduct
		var p models.Product
		err := rows.Scan(
			&op.ID, &op.OrderID, &op.ProductID, &op.Quantity, &op.Price,
			&p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return rows.Err()
}

func main() {
	dsn := "host=localhost user= password= dbname= port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	var customer models.Customer
	if err = db.
		QueryRow("SELECT id, name, email, created_at FROM customers LIMIT 1").
		Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.CreatedAt,
		); err != nil {
		log.Fatalf("failed to get customer: %v", err)
	}

	var product models.Product
	if err = db.
		QueryRow("SELECT id, name, description, price, stock, created_at FROM products LIMIT 1").
		Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
		); err != nil {
		log.Fatalf("failed to get product: %v", err)
	}

	var order models.Order
	if err = db.
		QueryRow("SELECT id, customer_id, date, total, created_at FROM orders LIMIT 1").
		Scan(
			&order.ID,
			&order.CustomerID,
			&order.Date,
			&order.Total,
			&order.CreatedAt,
		); err != nil {
		log.Fatalf("failed to get order: %v", err)
	}

	iterations := 10000

	operations := []utils.Operation{
		&CreateProduct{
			SQL: SQL{db},
		},
		&GetCustomerByID{
			ID:  customer.ID,
			SQL: SQL{db},
		},
		&UpdateProductByID{
			ID:    product.ID,
			Price: 89.99,
			SQL:   SQL{db},
		},
		&CreateOrderWithProductsByCustomerID{
			CustomerID: customer.ID,
			ProductID:  product.ID,
			SQL:        SQL{db},
		},
		&GetCustomerStatsByID{
			CustomerID: customer.ID,
			SQL:        SQL{db},
		},
		&GetProductSalesByLimit{
			Limit: 10,
			SQL:   SQL{db},
		},
		&GetOrderFullDetailsByID{
			OrderID: order.ID,
			SQL:     SQL{db},
		},
		&DeleteProductByName{
			SQL: SQL{db},
		},
	}

	utils.PrintResult(operations, iterations)
}
