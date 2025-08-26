package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// TagUser holds the schema definition for the user-tag entity.
type TagUser struct {
	ent.Schema
}

// Fields of the Tag.
func (TagUser) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tag_id").Optional(),
		field.Int("user_id").Optional(),
		field.Bool("is_read"),
		field.Bool("is_favorite"),
	}
}

// Edges of the Progress.
func (TagUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tag", Tag.Type).Ref("tag_user_details").Unique().Field("tag_id"),
		edge.From("user", User.Type).Ref("tag_user_details").Unique().Field("user_id"),
	}
}

func (TagUser) Indexes() []ent.Index {
	return []ent.Index{
		// Index for tag-user
		index.Fields("user_id", "tag_id").Unique(),
	}
}
