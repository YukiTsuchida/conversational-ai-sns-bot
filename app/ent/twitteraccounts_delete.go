// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/predicate"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/twitteraccounts"
)

// TwitterAccountsDelete is the builder for deleting a TwitterAccounts entity.
type TwitterAccountsDelete struct {
	config
	hooks    []Hook
	mutation *TwitterAccountsMutation
}

// Where appends a list predicates to the TwitterAccountsDelete builder.
func (tad *TwitterAccountsDelete) Where(ps ...predicate.TwitterAccounts) *TwitterAccountsDelete {
	tad.mutation.Where(ps...)
	return tad
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (tad *TwitterAccountsDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, tad.sqlExec, tad.mutation, tad.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (tad *TwitterAccountsDelete) ExecX(ctx context.Context) int {
	n, err := tad.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (tad *TwitterAccountsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(twitteraccounts.Table, sqlgraph.NewFieldSpec(twitteraccounts.FieldID, field.TypeInt))
	if ps := tad.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, tad.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	tad.mutation.done = true
	return affected, err
}

// TwitterAccountsDeleteOne is the builder for deleting a single TwitterAccounts entity.
type TwitterAccountsDeleteOne struct {
	tad *TwitterAccountsDelete
}

// Where appends a list predicates to the TwitterAccountsDelete builder.
func (tado *TwitterAccountsDeleteOne) Where(ps ...predicate.TwitterAccounts) *TwitterAccountsDeleteOne {
	tado.tad.mutation.Where(ps...)
	return tado
}

// Exec executes the deletion query.
func (tado *TwitterAccountsDeleteOne) Exec(ctx context.Context) error {
	n, err := tado.tad.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{twitteraccounts.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tado *TwitterAccountsDeleteOne) ExecX(ctx context.Context) {
	if err := tado.Exec(ctx); err != nil {
		panic(err)
	}
}
