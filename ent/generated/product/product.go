// Code generated by ent, DO NOT EDIT.

package product

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the product type in the database.
	Label = "product"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldPrice holds the string denoting the price field in the database.
	FieldPrice = "price"
	// FieldStock holds the string denoting the stock field in the database.
	FieldStock = "stock"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeOrderProducts holds the string denoting the order_products edge name in mutations.
	EdgeOrderProducts = "order_products"
	// Table holds the table name of the product in the database.
	Table = "products"
	// OrderProductsTable is the table that holds the order_products relation/edge.
	OrderProductsTable = "order_products"
	// OrderProductsInverseTable is the table name for the OrderProduct entity.
	// It exists in this package in order to avoid circular dependency with the "orderproduct" package.
	OrderProductsInverseTable = "order_products"
	// OrderProductsColumn is the table column denoting the order_products relation/edge.
	OrderProductsColumn = "product_id"
)

// Columns holds all SQL columns for product fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDescription,
	FieldPrice,
	FieldStock,
	FieldCreatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
)

// OrderOption defines the ordering options for the Product queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByPrice orders the results by the price field.
func ByPrice(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPrice, opts...).ToFunc()
}

// ByStock orders the results by the stock field.
func ByStock(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStock, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByOrderProductsCount orders the results by order_products count.
func ByOrderProductsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newOrderProductsStep(), opts...)
	}
}

// ByOrderProducts orders the results by order_products terms.
func ByOrderProducts(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOrderProductsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newOrderProductsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OrderProductsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, OrderProductsTable, OrderProductsColumn),
	)
}
