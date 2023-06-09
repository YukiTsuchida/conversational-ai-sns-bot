// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/conversations"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/twitteraccounts"
)

// TwitterAccounts is the model entity for the TwitterAccounts schema.
type TwitterAccounts struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// TwitterAccountID holds the value of the "twitter_account_id" field.
	TwitterAccountID string `json:"twitter_account_id,omitempty"`
	// AccessToken holds the value of the "access_token" field.
	AccessToken string `json:"access_token,omitempty"`
	// RefreshToken holds the value of the "refresh_token" field.
	RefreshToken string `json:"refresh_token,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the TwitterAccountsQuery when eager-loading is set.
	Edges                         TwitterAccountsEdges `json:"edges"`
	twitter_accounts_conversation *int
	selectValues                  sql.SelectValues
}

// TwitterAccountsEdges holds the relations/edges for other nodes in the graph.
type TwitterAccountsEdges struct {
	// Conversation holds the value of the conversation edge.
	Conversation *Conversations `json:"conversation,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ConversationOrErr returns the Conversation value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e TwitterAccountsEdges) ConversationOrErr() (*Conversations, error) {
	if e.loadedTypes[0] {
		if e.Conversation == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: conversations.Label}
		}
		return e.Conversation, nil
	}
	return nil, &NotLoadedError{edge: "conversation"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*TwitterAccounts) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case twitteraccounts.FieldID:
			values[i] = new(sql.NullInt64)
		case twitteraccounts.FieldTwitterAccountID, twitteraccounts.FieldAccessToken, twitteraccounts.FieldRefreshToken:
			values[i] = new(sql.NullString)
		case twitteraccounts.FieldCreatedAt, twitteraccounts.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case twitteraccounts.ForeignKeys[0]: // twitter_accounts_conversation
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the TwitterAccounts fields.
func (ta *TwitterAccounts) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case twitteraccounts.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ta.ID = int(value.Int64)
		case twitteraccounts.FieldTwitterAccountID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field twitter_account_id", values[i])
			} else if value.Valid {
				ta.TwitterAccountID = value.String
			}
		case twitteraccounts.FieldAccessToken:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field access_token", values[i])
			} else if value.Valid {
				ta.AccessToken = value.String
			}
		case twitteraccounts.FieldRefreshToken:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field refresh_token", values[i])
			} else if value.Valid {
				ta.RefreshToken = value.String
			}
		case twitteraccounts.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				ta.CreatedAt = value.Time
			}
		case twitteraccounts.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				ta.UpdatedAt = value.Time
			}
		case twitteraccounts.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field twitter_accounts_conversation", value)
			} else if value.Valid {
				ta.twitter_accounts_conversation = new(int)
				*ta.twitter_accounts_conversation = int(value.Int64)
			}
		default:
			ta.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the TwitterAccounts.
// This includes values selected through modifiers, order, etc.
func (ta *TwitterAccounts) Value(name string) (ent.Value, error) {
	return ta.selectValues.Get(name)
}

// QueryConversation queries the "conversation" edge of the TwitterAccounts entity.
func (ta *TwitterAccounts) QueryConversation() *ConversationsQuery {
	return NewTwitterAccountsClient(ta.config).QueryConversation(ta)
}

// Update returns a builder for updating this TwitterAccounts.
// Note that you need to call TwitterAccounts.Unwrap() before calling this method if this TwitterAccounts
// was returned from a transaction, and the transaction was committed or rolled back.
func (ta *TwitterAccounts) Update() *TwitterAccountsUpdateOne {
	return NewTwitterAccountsClient(ta.config).UpdateOne(ta)
}

// Unwrap unwraps the TwitterAccounts entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ta *TwitterAccounts) Unwrap() *TwitterAccounts {
	_tx, ok := ta.config.driver.(*txDriver)
	if !ok {
		panic("ent: TwitterAccounts is not a transactional entity")
	}
	ta.config.driver = _tx.drv
	return ta
}

// String implements the fmt.Stringer.
func (ta *TwitterAccounts) String() string {
	var builder strings.Builder
	builder.WriteString("TwitterAccounts(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ta.ID))
	builder.WriteString("twitter_account_id=")
	builder.WriteString(ta.TwitterAccountID)
	builder.WriteString(", ")
	builder.WriteString("access_token=")
	builder.WriteString(ta.AccessToken)
	builder.WriteString(", ")
	builder.WriteString("refresh_token=")
	builder.WriteString(ta.RefreshToken)
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(ta.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(ta.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// TwitterAccountsSlice is a parsable slice of TwitterAccounts.
type TwitterAccountsSlice []*TwitterAccounts
