package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Chatgpt35TurboConversationLog holds the schema definition for the Chatgpt35TurboConversationLog entity.
type Chatgpt35TurboConversationLog struct {
	ent.Schema
}

// Fields of the Chatgpt35TurboConversationLog.
func (Chatgpt35TurboConversationLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("message").NotEmpty(),
		field.String("purpose").Optional(),
		field.Enum("role").Values("system", "user", "assistant"),
		field.Time("created_at").
			Default(time.Now).
			Annotations(
				entsql.Default("CURRENT_TIMESTAMP"),
			),
	}
}

// Edges of the Chatgpt35TurboConversationLog.
func (Chatgpt35TurboConversationLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("conversation", Conversations.Type).Unique(),
	}
}
