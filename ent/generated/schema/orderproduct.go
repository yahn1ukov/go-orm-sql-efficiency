package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// OrderProduct holds the schema definition for the OrderProduct entity.
type OrderProduct struct {
	ent.Schema
}

// Fields of the OrderProduct.
func (OrderProduct) Fields() []ent.Field {
	return []ent.Field{
		field.Int("order_id"),
		field.Int("product_id"),
		field.Int("quantity"),
		field.Float("price"),
	}
}

// Edges of the OrderProduct.
func (OrderProduct) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).
			Ref("products").
			Field("order_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("product", Product.Type).
			Ref("order_products").
			Field("product_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
