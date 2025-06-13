package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yahn1ukov/go-orm-sql-efficiency/models"
	"github.com/yahn1ukov/go-orm-sql-efficiency/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GORM struct {
	db *gorm.DB
}

type CreateProduct struct {
	GORM
}

func (op *CreateProduct) Name() string {
	return "Create Product"
}

func (op *CreateProduct) Execute(iteration int) error {
	return op.db.
		Create(&models.Product{
			Name:        fmt.Sprintf("Product_%d", iteration),
			Description: "Test product",
			Price:       99.99,
			Stock:       100,
		}).
		Error
}

type GetCustomerByID struct {
	ID int
	GORM
}

func (op *GetCustomerByID) Name() string {
	return "Get Customer by ID"
}

func (op *GetCustomerByID) Execute(int) error {
	var customer models.Customer
	return op.db.
		First(&customer, op.ID).
		Error
}

type UpdateProductByID struct {
	ID    int
	Price float64
	GORM
}

func (op *UpdateProductByID) Name() string {
	return "Update Product Price by ID"
}

func (op *UpdateProductByID) Execute(int) error {
	return op.db.
		Model(&models.Product{}).
		Where("id = ?", op.ID).
		Update("price", op.Price).
		Error
}

type DeleteProductByName struct {
	GORM
}

func (op *DeleteProductByName) Name() string {
	return "Delete Product by Name"
}

func (op *DeleteProductByName) Execute(iteration int) error {
	return op.db.
		Where("name = ?", fmt.Sprintf("Product_%d", iteration)).
		Delete(&models.Product{}).
		Error
}

type CreateOrderWithProductsByCustomerID struct {
	CustomerID int
	ProductID  int
	GORM
}

func (op *CreateOrderWithProductsByCustomerID) Name() string {
	return "Create Order with Products By Customer ID (Transaction)"
}

func (op *CreateOrderWithProductsByCustomerID) Execute(int) error {
	return op.db.
		Transaction(func(tx *gorm.DB) error {
			return tx.
				Create(&models.Order{
					CustomerID: op.CustomerID,
					Date:       time.Now(),
					Total:      199.98,
					Products: []models.OrderProduct{
						{
							ProductID: op.ProductID,
							Quantity:  2,
							Price:     99.99,
						},
					},
				}).
				Error
		})
}

type GetCustomerStatsByID struct {
	CustomerID int
	GORM
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
		Model(&models.Order{}).
		Select("COUNT(*) as total_orders, COALESCE(SUM(total), 0) as total_spent").
		Where("customer_id = ?", op.CustomerID).
		Scan(&result).
		Error
}

type GetProductSalesByLimit struct {
	Limit int
	GORM
}

func (op *GetProductSalesByLimit) Name() string {
	return "Get Product Sales by Limit (Complex Join)"
}

func (op *GetProductSalesByLimit) Execute(int) error {
	type Result struct {
		ProductID   uint
		ProductName string
		TotalSales  int64
		Revenue     float64
	}

	var results []Result
	return op.db.
		Model(&models.OrderProduct{}).
		Select("product_id, products.name as product_name, SUM(quantity) as total_sales, SUM(order_products.price * quantity) as revenue").
		Joins("JOIN products ON products.id = order_products.product_id").
		Group("product_id, products.name").
		Order("revenue DESC").
		Limit(op.Limit).
		Scan(&results).
		Error
}

func main() {
	dsn := "host=localhost user= password= dbname= port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	var customer models.Customer
	if err = db.First(&customer).Error; err != nil {
		log.Fatalf("failed to get customer: %v", err)
	}

	var product models.Product
	if err = db.First(&product).Error; err != nil {
		log.Fatalf("failed to get product: %v", err)
	}

	iterations := 10000

	operations := []utils.Operation{
		&CreateProduct{
			GORM: GORM{db},
		},
		&GetCustomerByID{
			ID:   customer.ID,
			GORM: GORM{db},
		},
		&UpdateProductByID{
			ID:    product.ID,
			Price: 89.99,
			GORM:  GORM{db},
		},
		&CreateOrderWithProductsByCustomerID{
			CustomerID: customer.ID,
			ProductID:  product.ID,
			GORM:       GORM{db},
		},
		&GetCustomerStatsByID{
			CustomerID: customer.ID,
			GORM:       GORM{db},
		},
		&GetProductSalesByLimit{
			Limit: 10,
			GORM:  GORM{db},
		},
		&DeleteProductByName{
			GORM: GORM{db},
		},
	}

	utils.PrintResult(operations, iterations)
}
