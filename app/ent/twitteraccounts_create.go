// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/conversations"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/twitteraccounts"
)

// TwitterAccountsCreate is the builder for creating a TwitterAccounts entity.
type TwitterAccountsCreate struct {
	config
	mutation *TwitterAccountsMutation
	hooks    []Hook
}

// SetTwitterAccountID sets the "twitter_account_id" field.
func (tac *TwitterAccountsCreate) SetTwitterAccountID(s string) *TwitterAccountsCreate {
	tac.mutation.SetTwitterAccountID(s)
	return tac
}

// SetAccessToken sets the "access_token" field.
func (tac *TwitterAccountsCreate) SetAccessToken(s string) *TwitterAccountsCreate {
	tac.mutation.SetAccessToken(s)
	return tac
}

// SetRefreshToken sets the "refresh_token" field.
func (tac *TwitterAccountsCreate) SetRefreshToken(s string) *TwitterAccountsCreate {
	tac.mutation.SetRefreshToken(s)
	return tac
}

// SetCreatedAt sets the "created_at" field.
func (tac *TwitterAccountsCreate) SetCreatedAt(t time.Time) *TwitterAccountsCreate {
	tac.mutation.SetCreatedAt(t)
	return tac
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (tac *TwitterAccountsCreate) SetNillableCreatedAt(t *time.Time) *TwitterAccountsCreate {
	if t != nil {
		tac.SetCreatedAt(*t)
	}
	return tac
}

// SetUpdatedAt sets the "updated_at" field.
func (tac *TwitterAccountsCreate) SetUpdatedAt(t time.Time) *TwitterAccountsCreate {
	tac.mutation.SetUpdatedAt(t)
	return tac
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (tac *TwitterAccountsCreate) SetNillableUpdatedAt(t *time.Time) *TwitterAccountsCreate {
	if t != nil {
		tac.SetUpdatedAt(*t)
	}
	return tac
}

// SetConversationID sets the "conversation" edge to the Conversations entity by ID.
func (tac *TwitterAccountsCreate) SetConversationID(id int) *TwitterAccountsCreate {
	tac.mutation.SetConversationID(id)
	return tac
}

// SetNillableConversationID sets the "conversation" edge to the Conversations entity by ID if the given value is not nil.
func (tac *TwitterAccountsCreate) SetNillableConversationID(id *int) *TwitterAccountsCreate {
	if id != nil {
		tac = tac.SetConversationID(*id)
	}
	return tac
}

// SetConversation sets the "conversation" edge to the Conversations entity.
func (tac *TwitterAccountsCreate) SetConversation(c *Conversations) *TwitterAccountsCreate {
	return tac.SetConversationID(c.ID)
}

// Mutation returns the TwitterAccountsMutation object of the builder.
func (tac *TwitterAccountsCreate) Mutation() *TwitterAccountsMutation {
	return tac.mutation
}

// Save creates the TwitterAccounts in the database.
func (tac *TwitterAccountsCreate) Save(ctx context.Context) (*TwitterAccounts, error) {
	tac.defaults()
	return withHooks(ctx, tac.sqlSave, tac.mutation, tac.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tac *TwitterAccountsCreate) SaveX(ctx context.Context) *TwitterAccounts {
	v, err := tac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tac *TwitterAccountsCreate) Exec(ctx context.Context) error {
	_, err := tac.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tac *TwitterAccountsCreate) ExecX(ctx context.Context) {
	if err := tac.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tac *TwitterAccountsCreate) defaults() {
	if _, ok := tac.mutation.CreatedAt(); !ok {
		v := twitteraccounts.DefaultCreatedAt()
		tac.mutation.SetCreatedAt(v)
	}
	if _, ok := tac.mutation.UpdatedAt(); !ok {
		v := twitteraccounts.DefaultUpdatedAt()
		tac.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tac *TwitterAccountsCreate) check() error {
	if _, ok := tac.mutation.TwitterAccountID(); !ok {
		return &ValidationError{Name: "twitter_account_id", err: errors.New(`ent: missing required field "TwitterAccounts.twitter_account_id"`)}
	}
	if v, ok := tac.mutation.TwitterAccountID(); ok {
		if err := twitteraccounts.TwitterAccountIDValidator(v); err != nil {
			return &ValidationError{Name: "twitter_account_id", err: fmt.Errorf(`ent: validator failed for field "TwitterAccounts.twitter_account_id": %w`, err)}
		}
	}
	if _, ok := tac.mutation.AccessToken(); !ok {
		return &ValidationError{Name: "access_token", err: errors.New(`ent: missing required field "TwitterAccounts.access_token"`)}
	}
	if v, ok := tac.mutation.AccessToken(); ok {
		if err := twitteraccounts.AccessTokenValidator(v); err != nil {
			return &ValidationError{Name: "access_token", err: fmt.Errorf(`ent: validator failed for field "TwitterAccounts.access_token": %w`, err)}
		}
	}
	if _, ok := tac.mutation.RefreshToken(); !ok {
		return &ValidationError{Name: "refresh_token", err: errors.New(`ent: missing required field "TwitterAccounts.refresh_token"`)}
	}
	if v, ok := tac.mutation.RefreshToken(); ok {
		if err := twitteraccounts.RefreshTokenValidator(v); err != nil {
			return &ValidationError{Name: "refresh_token", err: fmt.Errorf(`ent: validator failed for field "TwitterAccounts.refresh_token": %w`, err)}
		}
	}
	if _, ok := tac.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "TwitterAccounts.created_at"`)}
	}
	if _, ok := tac.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "TwitterAccounts.updated_at"`)}
	}
	return nil
}

func (tac *TwitterAccountsCreate) sqlSave(ctx context.Context) (*TwitterAccounts, error) {
	if err := tac.check(); err != nil {
		return nil, err
	}
	_node, _spec := tac.createSpec()
	if err := sqlgraph.CreateNode(ctx, tac.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	tac.mutation.id = &_node.ID
	tac.mutation.done = true
	return _node, nil
}

func (tac *TwitterAccountsCreate) createSpec() (*TwitterAccounts, *sqlgraph.CreateSpec) {
	var (
		_node = &TwitterAccounts{config: tac.config}
		_spec = sqlgraph.NewCreateSpec(twitteraccounts.Table, sqlgraph.NewFieldSpec(twitteraccounts.FieldID, field.TypeInt))
	)
	if value, ok := tac.mutation.TwitterAccountID(); ok {
		_spec.SetField(twitteraccounts.FieldTwitterAccountID, field.TypeString, value)
		_node.TwitterAccountID = value
	}
	if value, ok := tac.mutation.AccessToken(); ok {
		_spec.SetField(twitteraccounts.FieldAccessToken, field.TypeString, value)
		_node.AccessToken = value
	}
	if value, ok := tac.mutation.RefreshToken(); ok {
		_spec.SetField(twitteraccounts.FieldRefreshToken, field.TypeString, value)
		_node.RefreshToken = value
	}
	if value, ok := tac.mutation.CreatedAt(); ok {
		_spec.SetField(twitteraccounts.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := tac.mutation.UpdatedAt(); ok {
		_spec.SetField(twitteraccounts.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if nodes := tac.mutation.ConversationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   twitteraccounts.ConversationTable,
			Columns: []string{twitteraccounts.ConversationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(conversations.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.twitter_accounts_conversation = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TwitterAccountsCreateBulk is the builder for creating many TwitterAccounts entities in bulk.
type TwitterAccountsCreateBulk struct {
	config
	builders []*TwitterAccountsCreate
}

// Save creates the TwitterAccounts entities in the database.
func (tacb *TwitterAccountsCreateBulk) Save(ctx context.Context) ([]*TwitterAccounts, error) {
	specs := make([]*sqlgraph.CreateSpec, len(tacb.builders))
	nodes := make([]*TwitterAccounts, len(tacb.builders))
	mutators := make([]Mutator, len(tacb.builders))
	for i := range tacb.builders {
		func(i int, root context.Context) {
			builder := tacb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TwitterAccountsMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, tacb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tacb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, tacb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tacb *TwitterAccountsCreateBulk) SaveX(ctx context.Context) []*TwitterAccounts {
	v, err := tacb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tacb *TwitterAccountsCreateBulk) Exec(ctx context.Context) error {
	_, err := tacb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tacb *TwitterAccountsCreateBulk) ExecX(ctx context.Context) {
	if err := tacb.Exec(ctx); err != nil {
		panic(err)
	}
}
