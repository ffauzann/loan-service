package model

import (
	"time"

	"github.com/ffauzann/loan-service/internal/constant"
)

type Loan struct {
	CommonModel

	BorrowerId      uint64             `json:"borrower_id" db:"borrower_id"`
	PrincipalAmount float64            `json:"principal_amount" db:"principal_amount"`
	InterestRate    float64            `json:"interest_rate" db:"interest_rate"`
	ROI             float64            `json:"roi" db:"roi"`
	AgreementLink   *string            `json:"agreement_link" db:"agreement_link"`
	State           constant.LoanState `json:"state" db:"state"`
	InvestedAmount  float64            `json:"invested_amount" db:"invested_amount"`
}

type LoanApproval struct {
	CommonModel

	LoanId         uint64    `json:"loan_id" db:"loan_id"`
	ValidatorId    uint64    `json:"validator_id" db:"validator_id"`
	PhotoProofLink string    `json:"photo_proof_link" db:"photo_proof_link"`
	ApprovalDate   time.Time `json:"approval_date" db:"approval_date"`
}

type LoanInvestment struct {
	CommonModel

	LoanId     uint64    `json:"loan_id" db:"loan_id"`
	InvestorId uint64    `json:"investor_id" db:"investor_id"`
	Amount     float64   `json:"amount" db:"amount"`
	InvestedAt time.Time `json:"invested_at" db:"invested_at"`
}

type LoanDisbursement struct {
	CommonModel

	LoanId              uint64    `json:"loan_id" db:"loan_id"`
	OfficerId           uint64    `json:"officer_id" db:"officer_id"`
	SignedAgreementLink string    `json:"signed_agreement_link" db:"signed_agreement_link"`
	DisbursementDate    time.Time `json:"disbursement_date" db:"disbursement_date"`
}

// -------------------- Create Loan --------------------

type CreateLoanRequest struct {
	BorrowerId      uint64  `json:"-"`                                                                      // comes from auth context
	PrincipalAmount float64 `json:"principal_amount" validate:"required,numeric,gte=1_000,lte=100_000_000"` // required
}

type CreateLoanResponse struct {
	LoanId    uint64    `json:"loan_id"`
	State     string    `json:"state"` // PROPOSED
	CreatedAt time.Time `json:"created_at"`
}

// -------------------- Approve Loan --------------------

type ApproveLoanRequest struct {
	LoanId         uint64  `json:"loan_id" validate:"required,gte=1"`                       // required
	ValidatorId    uint64  `json:"-"`                                                       // comes from auth context (Field Validator role)
	PhotoProofLink string  `json:"photo_proof_link" validate:"required,url"`                // required, link to photo proof of validation
	InterestRate   float64 `json:"interest_rate" validate:"required,numeric,gte=0,lte=100"` // % interest borrower pays
	ROI            float64 `json:"roi" validate:"required,numeric,gte=0,lte=100"`           // % return for investors
}

type ApproveLoanResponse struct {
	LoanId       uint64    `json:"loan_id"`
	State        string    `json:"state"` // APPROVED
	ApprovalDate time.Time `json:"approval_date"`
}

// -------------------- Invest in Loan --------------------

type InvestInLoanRequest struct {
	LoanId     uint64  `json:"loan_id"`
	InvestorId uint64  `json:"investor_id"` // comes from auth context
	Amount     float64 `json:"amount" validate:"required,numeric,gte=1000"`
}

type InvestInLoanResponse struct {
	LoanId         uint64  `json:"loan_id"`
	InvestedAmount float64 `json:"invested_amount"`
	State          string  `json:"state"` // FUNDING or INVESTED
}

// -------------------- Disburse Loan --------------------

type DisburseLoanRequest struct {
	LoanId              uint64 `json:"loan_id"`
	OfficerId           uint64 `json:"-"`                                             // comes from auth context (Admin role)
	SignedAgreementLink string `json:"signed_agreement_link" validate:"required,url"` // proof of signed contract
}

type DisburseLoanResponse struct {
	LoanId           uint64    `json:"loan_id"`
	State            string    `json:"state"` // DISBURSED
	DisbursementDate time.Time `json:"disbursement_date"`
}
