package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
)

// Conversations holds the schema definition for the Conversations entity.
type Conversations struct {
	ent.Schema
}

// Fields of the Conversations.
func (Conversations) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("ai_model").Values("chatgpt-3_5-turbo"),
		field.Enum("sns_type").Values("twitter"),
		field.Enum("cmd_version").Values("v0_1"),
		field.Bool("is_aborted").Default(false),
		field.Time("created_at").
			Default(time.Now).
			Annotations(
				entsql.Default("CURRENT_TIMESTAMP"),
			),
	}
}

// Edges of the Conversations.
func (Conversations) Edges() []ent.Edge {
	return nil
}
