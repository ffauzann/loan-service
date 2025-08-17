package client

// Known roleIds.
const (
	RoleIdSuperadmin uint8 = iota + 1
	RoleIdAdmin
	RoleIdUser
)

const (
	JwtIssuer = "example.com"
	JwksPath  = "/user/api/v1/r/utilities/.well-known/jwks.json"
)
