package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Tag holds the schema definition for the Tag entity.
type Serie struct {
	ent.Schema
}

// Fields of the Tag.
func (Serie) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty(),
		field.Bool("hidden").Default(false),
		field.Time("last_update").Default(time.Now),
	}
}

// Edges of the Tag.
func (Serie) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("items", Meta.Type).Ref("serie"),
		edge.From("favorite_of_user", User.Type).Ref("favorite_series"),
	}
}
