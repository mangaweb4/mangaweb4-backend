package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Tag holds the schema definition for the Tag entity.
type Tag struct {
	ent.Schema
}

// Fields of the Tag.
func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty(),
		field.Bool("favorite").Default(false).Deprecated("use 'favorite_of_user' edge instead."),
		field.Bool("hidden").Default(false),
		field.Time("last_update").Default(time.Time{}).Optional(),
	}
}

// Edges of the Tag.
func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("meta", Meta.Type).Ref("tags"),
		edge.From("favorite_of_user", User.Type).Ref("favorite_tags"),
		edge.To("tag_user_details", TagUser.Type),
	}
}
