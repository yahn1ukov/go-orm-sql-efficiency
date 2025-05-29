package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
	ent "github.com/yahn1ukov/go-orm-sql-efficiency/ent/generated"
	"github.com/yahn1ukov/go-orm-sql-efficiency/ent/generated/customer"
	"github.com/yahn1ukov/go-orm-sql-efficiency/ent/generated/order"
	"github.com/yahn1ukov/go-orm-sql-efficiency/ent/generated/product"
	"github.com/yahn1ukov/go-orm-sql-efficiency/utils"
)

type Ent struct {
	client *ent.Client
	db     *sql.DB
}

type CreateProduct struct {
	Ent
}

func (op *CreateProduct) Name() string {
	return "Create Product"
}

func (op *CreateProduct) Execute(iteration int) error {
	_, err := op.client.Product.
		Create().
		SetName(fmt.Sprintf("Product_%d", iteration)).
		SetDescription("Test product").
		SetPrice(99.99).
		SetStock(100).
		Save(context.Background())

	return err
}

type GetCustomerByID struct {
	ID int
	Ent
}

func (op *GetCustomerByID) Name() string {
	return "Get Customer by ID"
}

func (op *GetCustomerByID) Execute(int) error {
	_, err := op.client.Customer.
		Query().
		Where(customer.ID(op.ID)).
		Only(context.Background())

	return err
}

type UpdateProductByID struct {
	ID    int
	Price float64
	Ent
}

func (op *UpdateProductByID) Name() string {
	return "Update Product Price by ID"
}

func (op *UpdateProductByID) Execute(int) error {
	return op.client.Product.
		UpdateOneID(op.ID).
		SetPrice(op.Price).
		Exec(context.Background())
}

type DeleteProductByName struct {
	Ent
}

func (op *DeleteProductByName) Name() string {
	return "Delete Product by Name"
}

func (op *DeleteProductByName) Execute(iteration int) error {
	_, err := op.client.Product.Delete().
		Where(product.Name(fmt.Sprintf("Product_%d", iteration))).
		Exec(context.Background())

	return err
}

type CreateOrderWithProductsByCustomerID struct {
	CustomerID int
	ProductID  int
	Ent
}

func (op *CreateOrderWithProductsByCustomerID) Name() string {
	return "Create Order with Products By Customer ID (Transaction)"
}

func (op *CreateOrderWithProductsByCustomerID) Execute(int) error {
	tx, err := op.client.Tx(context.Background())
	if err != nil {
		return err
	}

	o, err := tx.Order.
		Create().
		SetCustomerID(op.CustomerID).
		SetDate(time.Now()).
		SetTotal(199.98).
		Save(context.Background())
	if err != nil {
		return rollback(tx, err)
	}

	_, err = tx.OrderProduct.
		Create().
		SetOrderID(o.ID).
		SetProductID(op.ProductID).
		SetQuantity(2).
		SetPrice(99.99).
		Save(context.Background())
	if err != nil {
		return rollback(tx, err)
	}

	return tx.Commit()
}

func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}

	return err
}

type GetCustomerStatsByID struct {
	CustomerID int
	Ent
}

func (op *GetCustomerStatsByID) Name() string {
	return "Get Customer Stats by ID (Aggregation)"
}

func (op *GetCustomerStatsByID) Execute(int) error {
	rows, err := op.db.
		QueryContext(
			context.Background(),
			"SELECT COUNT(*) AS total_orders, COALESCE(SUM(total), 0) AS total_spent FROM orders WHERE customer_id = $1",
			op.CustomerID,
		)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		var totalOrders int64
		var totalSpent float64
		if err := rows.Scan(&totalOrders, &totalSpent); err != nil {
			return err
		}
	}

	return rows.Err()
}

type GetProductSalesByLimit struct {
	Limit int
	Ent
}

func (op *GetProductSalesByLimit) Name() string {
	return "Get Product Sales by Limit (Complex Join)"
}

func (op *GetProductSalesByLimit) Execute(int) error {
	rows, err := op.db.
		QueryContext(
			context.Background(),
			`SELECT 
            p.id, 
            p.name, 
            SUM(op.quantity) AS total_sales, 
            SUM(op.price * op.quantity) AS revenue
        FROM order_products op
        JOIN products p ON p.id = op.product_id
        GROUP BY p.id, p.name
        ORDER BY revenue DESC
        LIMIT $1`,
			op.Limit,
		)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var productID int
		var productName string
		var totalSales int64
		var revenue float64
		if err := rows.Scan(
			&productID,
			&productName,
			&totalSales,
			&revenue,
		); err != nil {
			return err
		}
	}

	return rows.Err()
}

type GetOrderFullDetailsByID struct {
	OrderID int
	Ent
}

func (op *GetOrderFullDetailsByID) Name() string {
	return "Get Order Full Details by ID (Nested Preload)"
}

func (op *GetOrderFullDetailsByID) Execute(int) error {
	_, err := op.client.Order.
		Query().
		Where(order.ID(op.OrderID)).
		WithProducts(func(q *ent.OrderProductQuery) {
			q.WithProduct()
		}).
		Only(context.Background())

	return err
}

func main() {
	dsn := "host=localhost user= password= dbname= port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)

	client := ent.NewClient(ent.Driver(drv))
	defer client.Close()

	ctx := context.Background()

	customerEntity, err := client.Customer.
		Query().
		First(ctx)
	if err != nil {
		log.Fatalf("failed to get customer: %v", err)
	}

	productEntity, err := client.Product.
		Query().
		First(ctx)
	if err != nil {
		log.Fatalf("failed to get product: %v", err)
	}

	orderEntity, err := client.Order.
		Query().
		First(ctx)
	if err != nil {
		log.Fatalf("failed to get order: %v", err)
	}

	iterations := 10000

	clients := Ent{
		client: client,
		db:     db,
	}

	operations := []utils.Operation{
		&CreateProduct{
			Ent: clients,
		},
		&GetCustomerByID{
			ID:  customerEntity.ID,
			Ent: clients,
		},
		&UpdateProductByID{
			ID:    productEntity.ID,
			Price: 89.99,
			Ent:   clients,
		},
		&CreateOrderWithProductsByCustomerID{
			CustomerID: customerEntity.ID,
			ProductID:  productEntity.ID,
			Ent:        clients,
		},
		&GetCustomerStatsByID{
			CustomerID: customerEntity.ID,
			Ent:        clients,
		},
		&GetProductSalesByLimit{
			Limit: 10,
			Ent:   clients,
		},
		&GetOrderFullDetailsByID{
			OrderID: orderEntity.ID,
			Ent:     clients,
		},
		&DeleteProductByName{
			Ent: clients,
		},
	}

	utils.PrintResult(operations, iterations)
}
