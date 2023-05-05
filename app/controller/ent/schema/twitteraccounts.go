package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TwitterAccounts holds the schema definition for the TwitterAccounts entity.
type TwitterAccounts struct {
	ent.Schema
}

// Fields of the TwitterAccounts.
func (TwitterAccounts) Fields() []ent.Field {
	return []ent.Field{
		field.String("twitter_account_id").NotEmpty().Unique(),
		field.String("bearer_token").NotEmpty(),
		field.Time("created_at").
			Default(time.Now).
			Annotations(
				entsql.Default("CURRENT_TIMESTAMP"),
			),
		field.Time("updated_at").
			Default(time.Now).
			Annotations(
				entsql.Default("CURRENT_TIMESTAMP"),
			),
	}
}

// Edges of the TwitterAccounts.
func (TwitterAccounts) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("conversation", Conversations.Type).Unique(),
	}
}
