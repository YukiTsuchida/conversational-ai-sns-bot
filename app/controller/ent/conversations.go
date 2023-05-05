// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent/conversations"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent/twitteraccounts"
)

// Conversations is the model entity for the Conversations schema.
type Conversations struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// AiModel holds the value of the "ai_model" field.
	AiModel conversations.AiModel `json:"ai_model,omitempty"`
	// SnsType holds the value of the "sns_type" field.
	SnsType conversations.SnsType `json:"sns_type,omitempty"`
	// CmdVersion holds the value of the "cmd_version" field.
	CmdVersion conversations.CmdVersion `json:"cmd_version,omitempty"`
	// IsAborted holds the value of the "is_aborted" field.
	IsAborted bool `json:"is_aborted,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ConversationsQuery when eager-loading is set.
	Edges                         ConversationsEdges `json:"edges"`
	twitter_accounts_conversation *int
	selectValues                  sql.SelectValues
}

// ConversationsEdges holds the relations/edges for other nodes in the graph.
type ConversationsEdges struct {
	// TwitterAccount holds the value of the twitter_account edge.
	TwitterAccount *TwitterAccounts `json:"twitter_account,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// TwitterAccountOrErr returns the TwitterAccount value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ConversationsEdges) TwitterAccountOrErr() (*TwitterAccounts, error) {
	if e.loadedTypes[0] {
		if e.TwitterAccount == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: twitteraccounts.Label}
		}
		return e.TwitterAccount, nil
	}
	return nil, &NotLoadedError{edge: "twitter_account"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Conversations) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case conversations.FieldIsAborted:
			values[i] = new(sql.NullBool)
		case conversations.FieldID:
			values[i] = new(sql.NullInt64)
		case conversations.FieldAiModel, conversations.FieldSnsType, conversations.FieldCmdVersion:
			values[i] = new(sql.NullString)
		case conversations.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		case conversations.ForeignKeys[0]: // twitter_accounts_conversation
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Conversations fields.
func (c *Conversations) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case conversations.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = int(value.Int64)
		case conversations.FieldAiModel:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field ai_model", values[i])
			} else if value.Valid {
				c.AiModel = conversations.AiModel(value.String)
			}
		case conversations.FieldSnsType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field sns_type", values[i])
			} else if value.Valid {
				c.SnsType = conversations.SnsType(value.String)
			}
		case conversations.FieldCmdVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field cmd_version", values[i])
			} else if value.Valid {
				c.CmdVersion = conversations.CmdVersion(value.String)
			}
		case conversations.FieldIsAborted:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_aborted", values[i])
			} else if value.Valid {
				c.IsAborted = value.Bool
			}
		case conversations.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				c.CreatedAt = value.Time
			}
		case conversations.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field twitter_accounts_conversation", value)
			} else if value.Valid {
				c.twitter_accounts_conversation = new(int)
				*c.twitter_accounts_conversation = int(value.Int64)
			}
		default:
			c.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Conversations.
// This includes values selected through modifiers, order, etc.
func (c *Conversations) Value(name string) (ent.Value, error) {
	return c.selectValues.Get(name)
}

// QueryTwitterAccount queries the "twitter_account" edge of the Conversations entity.
func (c *Conversations) QueryTwitterAccount() *TwitterAccountsQuery {
	return NewConversationsClient(c.config).QueryTwitterAccount(c)
}

// Update returns a builder for updating this Conversations.
// Note that you need to call Conversations.Unwrap() before calling this method if this Conversations
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Conversations) Update() *ConversationsUpdateOne {
	return NewConversationsClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Conversations entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Conversations) Unwrap() *Conversations {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Conversations is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Conversations) String() string {
	var builder strings.Builder
	builder.WriteString("Conversations(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("ai_model=")
	builder.WriteString(fmt.Sprintf("%v", c.AiModel))
	builder.WriteString(", ")
	builder.WriteString("sns_type=")
	builder.WriteString(fmt.Sprintf("%v", c.SnsType))
	builder.WriteString(", ")
	builder.WriteString("cmd_version=")
	builder.WriteString(fmt.Sprintf("%v", c.CmdVersion))
	builder.WriteString(", ")
	builder.WriteString("is_aborted=")
	builder.WriteString(fmt.Sprintf("%v", c.IsAborted))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(c.CreatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// ConversationsSlice is a parsable slice of Conversations.
type ConversationsSlice []*Conversations
