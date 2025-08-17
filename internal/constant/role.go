package constant

// Known roleIds.
const (
	RoleIdSuperadmin uint8 = iota + 1
	RoleIdAdmin
	RoleIdFieldValidator
	RoleIdInvestor
	RoleIdBorrower
)

// RBAC.
var (
	// Registration.
	AllowedRolesRegister    = []uint8{RoleIdSuperadmin, RoleIdAdmin}
	AllowedRolesRegisterMap = map[uint8][]uint8{
		RoleIdSuperadmin: {RoleIdSuperadmin, RoleIdAdmin},
		RoleIdAdmin:      {RoleIdAdmin},
	}

	// Propose loan.
	AllowedRolesProposeLoan = []uint8{RoleIdSuperadmin, RoleIdAdmin, RoleIdBorrower}

	// Approve loan.
	AllowedRolesApproveLoan = []uint8{RoleIdSuperadmin, RoleIdAdmin, RoleIdFieldValidator}

	// Invest loan.
	AllowedRolesInvestLoan = []uint8{RoleIdSuperadmin, RoleIdAdmin, RoleIdInvestor}

	// Disburse loan.
	AllowedRolesDisburseLoan = []uint8{RoleIdSuperadmin, RoleIdAdmin}
)
