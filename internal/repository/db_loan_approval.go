package repository

import (
	"context"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/model"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/util"
	"github.com/jmoiron/sqlx"
)

func (r *dbRepository) ApproveLoan(ctx context.Context, approval *model.LoanApproval, tx *sqlx.Tx) (err error) {
	if tx == nil { // End tx as soon as this method finishes if tx was not provided.
		defer func() { r.EndTx(ctx, tx, err) }()
	}

	tx, err = r.useOrInitTx(ctx, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	query := `
	INSERT INTO loan_approval (
		loan_id,
		validator_id,
		photo_proof_link,
		approval_date,
		created_by
	) VALUES (
		:loan_id,
		:validator_id,
		:photo_proof_link,
		:approval_date,
		:created_by
	)
	RETURNING id
	`

	insertApprovalQuery, args, err := tx.BindNamed(query, approval)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	err = tx.QueryRowxContext(ctx, insertApprovalQuery, args...).Scan(&approval.Id)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	return
}
