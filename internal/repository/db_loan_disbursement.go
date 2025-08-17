package repository

import (
	"context"

	"github.com/ffauzann/loan-service/internal/model"
	"github.com/ffauzann/loan-service/internal/util"
	"github.com/jmoiron/sqlx"
)

func (r *dbRepository) CreateLoanDisbursement(ctx context.Context, disbursement *model.LoanDisbursement, tx *sqlx.Tx) (err error) {
	if tx == nil { // End tx as soon as this method finishes if tx was not provided.
		defer func() { r.EndTx(ctx, tx, err) }()
	}

	tx, err = r.useOrInitTx(ctx, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	query := `
	INSERT INTO loan_disbursement (
		loan_id,
		officer_id,
		signed_agreement_link,
		disbursement_date,
		created_by
	) VALUES (
		:loan_id,
		:officer_id,
		:signed_agreement_link,
		:disbursement_date,
		:created_by
	)
	RETURNING id
	`

	insertApprovalQuery, args, err := tx.BindNamed(query, disbursement)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	err = tx.QueryRowxContext(ctx, insertApprovalQuery, args...).Scan(&disbursement.Id)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	return
}
