package repository

import (
	"context"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/model"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/util"
	"github.com/jmoiron/sqlx"
)

func (r *dbRepository) CreateLoanInvestment(ctx context.Context, investment *model.LoanInvestment, tx *sqlx.Tx) (err error) {
	if tx == nil { // End tx as soon as this method finishes if tx was not provided.
		defer func() { r.EndTx(ctx, tx, err) }()
	}

	tx, err = r.useOrInitTx(ctx, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	query := `
	INSERT INTO loan_investment (
		loan_id,
		investor_id,
		amount,
		created_by
	) VALUES (
		:loan_id,
		:investor_id,
		:amount,
		:created_by
	)
	RETURNING id
	`
	insertInvestmentQuery, args, err := tx.BindNamed(query, investment)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	err = tx.QueryRowxContext(ctx, insertInvestmentQuery, args...).Scan(&investment.Id)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	return
}

// GetInvesmentsByLoanId returns all investments for a given loan ID.
func (r *dbRepository) GetInvestmentsByLoanId(ctx context.Context, loanId uint64, tx *sqlx.Tx) (investments []*model.LoanInvestment, err error) {
	if tx == nil { // End tx as soon as this method finishes if tx was not provided.
		defer func() { r.EndTx(ctx, tx, err) }()
	}

	tx, err = r.useOrInitTx(ctx, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	query := `
	SELECT *
	FROM loan_investment
	WHERE loan_id = $1
	`

	if err = tx.SelectContext(ctx, &investments, query, loanId); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	return
}
