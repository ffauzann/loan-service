package repository

import (
	"context"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/model"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/util"
	"github.com/jmoiron/sqlx"
)

// CreateLoan inserts a new loan record into the database.
func (r *dbRepository) CreateLoan(ctx context.Context, loan *model.Loan, tx *sqlx.Tx) (err error) {
	if tx == nil { // End tx as soon as this method finishes if tx was not provided.
		defer func() { r.EndTx(ctx, tx, err) }()
	}

	tx, err = r.useOrInitTx(ctx, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	query := `
	INSERT INTO loan (
		borrower_id, 
		principal_amount, 
		interest_rate, 
		roi, 
		agreement_link, 
		state, 
		invested_amount, 
		created_by,
		updated_by
	)
	VALUES (
		:borrower_id, 
		:principal_amount, 
		:interest_rate, 
		:roi, 
		:agreement_link, 
		:state, 
		:invested_amount, 
		:created_by,
		:updated_by
	)
	RETURNING id
	`

	query, args, err := tx.BindNamed(query, loan)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	err = tx.QueryRowxContext(ctx, query, args...).Scan(&loan.Id)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	return
}

func (r *dbRepository) GetLoanById(ctx context.Context, loanID uint64, tx *sqlx.Tx) (loan *model.Loan, err error) {
	loan = &model.Loan{}
	// Use provided transaction or init one if nil
	tx, err = r.useOrInitTx(ctx, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return nil, err
	}

	query := `
	SELECT *
	FROM loan
	WHERE id = $1 AND deleted_at IS NULL
	`

	err = tx.GetContext(ctx, loan, query, loanID)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	return
}

// UpdateLoan updates an existing loan record in the database.
func (r *dbRepository) UpdateLoan(ctx context.Context, loan *model.Loan, tx *sqlx.Tx) (err error) {
	if tx == nil { // End tx as soon as this method finishes if tx was not provided.
		defer func() { r.EndTx(ctx, tx, err) }()
	}

	tx, err = r.useOrInitTx(ctx, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	query := `
	UPDATE loan
	SET
		borrower_id       = :borrower_id,
		principal_amount  = :principal_amount,
		interest_rate     = :interest_rate,
		roi               = :roi,
		agreement_link    = :agreement_link,
		state             = :state,
		invested_amount   = :invested_amount,
		updated_at        = NOW(),
		updated_by        = :updated_by
	WHERE id = :id
	`

	query, args, err := tx.BindNamed(query, loan)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	return
}
