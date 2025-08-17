package service

import (
	"context"
	"database/sql"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/constant"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/model"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/util"
)

func (s *service) CreateLoan(ctx context.Context, req *model.CreateLoanRequest) (res *model.CreateLoanResponse, err error) {
	// Prepare loan model
	loan := &model.Loan{
		BorrowerId:      req.BorrowerId,
		PrincipalAmount: req.PrincipalAmount,
		State:           constant.LoanStateProposed,
		CommonModel: model.CommonModel{
			CreatedAt: now(),
			CreatedBy: sql.NullInt64{Int64: int64(req.BorrowerId), Valid: true},
			UpdatedAt: sql.NullTime{Time: now(), Valid: true},
			UpdatedBy: sql.NullInt64{Int64: int64(req.BorrowerId), Valid: true},
		},
	}

	// Create loan.
	err = s.repository.db.CreateLoan(ctx, loan, nil)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = &model.CreateLoanResponse{
		LoanId:    loan.Id,
		State:     string(constant.LoanStateProposed),
		CreatedAt: loan.CreatedAt,
	}

	return
}

func (s *service) ApproveLoan(ctx context.Context, req *model.ApproveLoanRequest) (res *model.ApproveLoanResponse, err error) {
	// Get loan by ID.
	loan, err := s.repository.db.GetLoanById(ctx, req.LoanId, nil)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Validate loan state.
	if loan.State != constant.LoanStateProposed {
		err = constant.ErrLoanNotProposed
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	// Prepare loan model
	loanApproval := &model.LoanApproval{
		ValidatorId:    req.ValidatorId,
		LoanId:         req.LoanId,
		PhotoProofLink: req.PhotoProofLink,
		ApprovalDate:   now(),
		CommonModel: model.CommonModel{
			CreatedAt: now(),
			CreatedBy: sql.NullInt64{Int64: int64(req.ValidatorId), Valid: true},
			UpdatedAt: sql.NullTime{Time: now(), Valid: true},
			UpdatedBy: sql.NullInt64{Int64: int64(req.ValidatorId), Valid: true},
		},
	}

	// Set loan state to approved.
	loan.State = constant.LoanStateApproved
	loan.UpdatedBy = sql.NullInt64{Int64: int64(req.ValidatorId), Valid: true}

	// Begin tx.
	tx, err := s.repository.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}
	defer func() { s.repository.db.EndTx(ctx, tx, err) }()

	// Update loan state to approved.
	err = s.repository.db.UpdateLoan(ctx, loan, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Create loan approval.
	err = s.repository.db.ApproveLoan(ctx, loanApproval, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = &model.ApproveLoanResponse{
		LoanId:       loanApproval.LoanId,
		State:        string(constant.LoanStateApproved),
		ApprovalDate: loanApproval.ApprovalDate,
	}

	return
}

// InvestInLoan handles the investment in a loan proposal.
func (s *service) InvestInLoan(ctx context.Context, req *model.InvestInLoanRequest) (res *model.InvestInLoanResponse, err error) {
	// Begin tx.
	tx, err := s.repository.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}
	defer func() { s.repository.db.EndTx(ctx, tx, err) }()

	// Get loan by ID.
	loan, err := s.repository.db.GetLoanById(ctx, req.LoanId, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Validate loan state.
	if loan.State != constant.LoanStateApproved && loan.State != constant.LoanStateFunding {
		err = constant.ErrLoanNotApproved
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	// Validate investment amount.
	if req.Amount < constant.MinInvestmentAmount || req.Amount > constant.MaxInvestmentAmount {
		err = constant.ErrInvestmentAmountOutOfRange
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	// Validate funding progress.
	finalInvestedAmount := loan.InvestedAmount + req.Amount
	if finalInvestedAmount > loan.PrincipalAmount {
		err = constant.ErrInvestmentAmountOutOfRange
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	// Validate if the loan is fully funded.
	if loan.InvestedAmount+req.Amount == loan.PrincipalAmount {
		loan.State = constant.LoanStateInvested
	} else {
		loan.State = constant.LoanStateFunding
	}

	// Prepare investment model.
	investment := &model.LoanInvestment{
		LoanId:     req.LoanId,
		InvestorId: req.InvestorId,
		Amount:     req.Amount,
		CommonModel: model.CommonModel{
			CreatedAt: now(),
			CreatedBy: sql.NullInt64{Int64: int64(req.InvestorId), Valid: true},
			UpdatedAt: sql.NullTime{Time: now(), Valid: true},
			UpdatedBy: sql.NullInt64{Int64: int64(req.InvestorId), Valid: true},
		},
	}

	// Save investment record.
	err = s.repository.db.CreateLoanInvestment(ctx, investment, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Update loan with new invested amount and state.
	loan.InvestedAmount = finalInvestedAmount
	loan.UpdatedBy = sql.NullInt64{Int64: int64(req.InvestorId), Valid: true}
	err = s.repository.db.UpdateLoan(ctx, loan, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	if loan.State == constant.LoanStateInvested {
		// Notify all investors about the loan being fully funded.
		go s.notifyLoanFullyFunded(context.Background(), loan.Id)
	}

	// Construct response.
	res = &model.InvestInLoanResponse{
		LoanId:         req.LoanId,
		InvestedAmount: req.Amount,
		State:          string(loan.State),
	}

	return
}

// DisburseLoan handles the disbursement of a loan.
func (s *service) DisburseLoan(ctx context.Context, req *model.DisburseLoanRequest) (res *model.DisburseLoanResponse, err error) {
	// Begin tx.
	tx, err := s.repository.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}
	defer func() { s.repository.db.EndTx(ctx, tx, err) }()

	// Get loan by ID.
	loan, err := s.repository.db.GetLoanById(ctx, req.LoanId, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Validate loan state.
	if loan.State != constant.LoanStateInvested {
		err = constant.ErrLoanNotFullyInvested
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	// Prepare investment model.
	disbursement := &model.LoanDisbursement{
		LoanId:              req.LoanId,
		OfficerId:           req.OfficerId,
		SignedAgreementLink: req.SignedAgreementLink,
		DisbursementDate:    now(),
		CommonModel: model.CommonModel{
			CreatedAt: now(),
			CreatedBy: sql.NullInt64{Int64: int64(req.OfficerId), Valid: true},
			UpdatedAt: sql.NullTime{Time: now(), Valid: true},
			UpdatedBy: sql.NullInt64{Int64: int64(req.OfficerId), Valid: true},
		},
	}

	// Save investment record.
	err = s.repository.db.CreateLoanDisbursement(ctx, disbursement, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Update loan with new invested amount and state.
	loan.State = constant.LoanStateDisbursed
	loan.UpdatedBy = sql.NullInt64{Int64: int64(req.OfficerId), Valid: true}
	err = s.repository.db.UpdateLoan(ctx, loan, tx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = &model.DisburseLoanResponse{
		LoanId:           req.LoanId,
		DisbursementDate: disbursement.DisbursementDate,
		State:            string(loan.State),
	}

	return
}
