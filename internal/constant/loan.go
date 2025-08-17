package constant

// LoanState represents the state of a loan in the system.
type LoanState string

const (
	LoanStateProposed  LoanState = "PROPOSED"
	LoanStateApproved  LoanState = "APPROVED"
	LoanStateFunding   LoanState = "FUNDING"
	LoanStateInvested  LoanState = "INVESTED"
	LoanStateDisbursed LoanState = "DISBURSED"
)

const (
	MinInvestmentAmount float64 = 1_000.0       // Minimum investment amount in the system.
	MaxInvestmentAmount float64 = 100_000_000.0 // Maximum investment amount in the system.
)
