package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the Tag entity.
type User struct {
	ent.Schema
}

// Fields of the Tag.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").Unique().NotEmpty(),
		field.Bool("active").Default(true),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("favorite_items", Meta.Type),
		edge.To("favorite_tags", Tag.Type),
		edge.To("histories", History.Type),
		edge.To("progress", Progress.Type),
		edge.To("tag_user_details", TagUser.Type),
	}
}
